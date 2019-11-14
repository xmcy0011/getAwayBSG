package entrance

import (
	"encoding/json"
	"github.com/getAwayBSG/configs"
	"github.com/getAwayBSG/db"
	"github.com/getAwayBSG/logger"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	cachemongo "github.com/zolamk/colly-mongo-storage/colly/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	TotalPage int
	CurPage   int
}

func getAeraUrl(cityUrl string, areaListChan chan []string) {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.OnHTML("body", func(element *colly.HTMLElement) {
		element.ForEachWithBreak(".position div div", func(i int, element *colly.HTMLElement) bool {
			u, err := url.Parse(cityUrl)
			if err != nil {
				panic(err)
			}
			areaArr := make([]string, 0)

			element.ForEach("a", func(i int, element *colly.HTMLElement) {
				goUrl := element.Attr("href")
				rootUrl := u.Scheme + "://" + u.Host
				areaArr = append(areaArr, rootUrl+goUrl)
			})

			logger.Sugar.Infof("%s,抓取地区共:%d", cityUrl, len(areaArr))
			areaListChan <- areaArr
			return false
		})
	})

	err := c.Visit(cityUrl)
	if err != nil {
		logger.Sugar.Fatalf("%s:%s", err.Error(), cityUrl)
	}
}

func crawlerOneCity(cityUrl string, index int, total int) {
	c := colly.NewCollector()
	configInfo := configs.Config()

	if configInfo["crawlDelay"] != nil {
		delay, _ := configInfo["crawlDelay"].(json.Number).Int64()
		if delay > 0 {
			c.Limit(&colly.LimitRule{
				DomainGlob: "*",
				Delay:      time.Duration(delay) * time.Second,
			})
		}
	}

	if configInfo["proxyList"] != nil && len(configInfo["proxyList"].([]interface{})) > 0 {
		var proxyList []string
		for _, v := range configInfo["proxyList"].([]interface{}) {
			proxyList = append(proxyList, v.(string))
		}

		if configInfo["proxyList"] != nil {
			rp, err := proxy.RoundRobinProxySwitcher(proxyList...)
			if err != nil {
				logger.Sugar.Error(err)
			}
			c.SetProxyFunc(rp)
		}
	}
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	storage := &cachemongo.Storage{
		Database: "colly",
		URI:      configInfo["dburl"].(string) + "/colly",
	}
	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}

	areaListChan := make(chan []string, 1)
	getAeraUrl(cityUrl, areaListChan)
	areaArr := <-areaListChan

	var page Page
	var areaName string

	arr := strings.Split(cityUrl, ".")
	cityName := "unknown_city"
	if len(arr) > 1 {
		cityName = arr[0][8:]
	}

	// 挨个爬各个地区的房源
	for i := range areaArr {
		c.OnHTML("body", func(element *colly.HTMLElement) {
			// 一个地区的总数
			element.ForEach(".total", func(i int, element *colly.HTMLElement) {
				areaName = element.Text
			})
			// 获取总页数
			element.ForEach(".page-box", func(i int, element *colly.HTMLElement) {
				var tempPage Page
				err := json.Unmarshal([]byte(element.ChildAttr(".house-lst-page-box", "page-data")), &tempPage)
				if err == nil {
					page = tempPage
					logger.Sugar.Infof("[1/2][%d/%d][%s] totalCount=%s,totalPage=%d,curPage=%d",
						index+1, total, cityName, areaName, page.TotalPage, page.CurPage)
				}
			})

			// 获取一页的数据
			curCount := 0
			element.ForEach(".LOGCLICKDATA", func(i int, e *colly.HTMLElement) {
				link := e.ChildAttr("a", "href")

				title := e.ChildText("a:first-child")
				if title == "" {
					return
				}
				curCount++

				price := e.ChildText(".totalPrice")
				price = strings.Replace(price, "万", "0000", 1)
				iPrice, err := strconv.Atoi(price)
				if err != nil {
					iPrice = 0
				}

				unitPrice := e.ChildAttr(".unitPrice", "data-price")
				iUnitPrice, err := strconv.Atoi(unitPrice)
				if err != nil {
					iUnitPrice = 0
				}

				logger.Sugar.Infof("[1/2][%d/%d][%d/%d][%d] %s,%s,%s,总价：%d 万元，每平米：%d", index+1, total,
					page.CurPage, page.TotalPage, curCount, cityName, areaName, title, iPrice, iUnitPrice)

				db.Add(bson.M{"zq_detail_status": 0, "Title": title, "TotalePrice": iPrice, "UnitPrice": iUnitPrice, "Link": link, "listCrawlTime": time.Now()})
			})

			// 下一页
			element.ForEach(".page-box", func(i int, element *colly.HTMLElement) {
				var tempPage Page
				err := json.Unmarshal([]byte(element.ChildAttr(".house-lst-page-box", "page-data")), &tempPage)
				if err == nil {
					if page.CurPage < page.TotalPage {
						re, _ := regexp.Compile("pg\\d+/*")
						nextPageUrl := re.ReplaceAllString(element.Request.URL.String(), "")
						nextPageUrl = nextPageUrl + "pg" + strconv.Itoa(page.CurPage+1)
						err = c.Visit(nextPageUrl)
						if err != nil {
							logger.Sugar.Info(err)
						}
					}
				}
			})
		})

		err := c.Visit(areaArr[i])
		if err != nil {
			logger.Sugar.Debugf("%s:%s", err.Error(), cityUrl)
		}
	}
}

func listCrawler() {
	confInfo := configs.Config()
	cityList := confInfo["cityList"].([]interface{})
	count := len(cityList)
	for i := 0; i < count; i++ {
		cityName := cityList[i].(string)
		logger.Sugar.Infof("[1/2][%d/%d] 抓取城市：%s", i+1, count, cityName)
		crawlerOneCity(cityName, i, count)
	}
}

func crawlDetail() (sucnum int) {
	sucnum = 0
	c := colly.NewCollector()
	configInfo := configs.Config()

	//设置延时
	if configInfo["crawlDelay"] != nil {
		delay, _ := configInfo["crawlDelay"].(json.Number).Int64()
		if delay > 0 {
			c.Limit(&colly.LimitRule{
				DomainGlob: "*",
				Delay:      time.Duration(delay) * time.Second,
			})
		}
	}

	//设置代理
	if configInfo["proxyList"] != nil && len(configInfo["proxyList"].([]interface{})) > 0 {
		var proxyList []string
		for _, v := range configInfo["proxyList"].([]interface{}) {
			proxyList = append(proxyList, v.(string))
		}

		if configInfo["proxyList"] != nil {
			rp, err := proxy.RoundRobinProxySwitcher(proxyList...)
			if err != nil {
				logger.Sugar.Error(err)
			}
			c.SetProxyFunc(rp)
		}
	}

	//随机UA
	extensions.RandomUserAgent(c)
	//自动referer
	extensions.Referer(c)
	//设置MongoDB存储状态信息
	storage := &cachemongo.Storage{
		Database: "colly",
		URI:      configInfo["dburl"].(string) + "/colly",
	}
	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}
	c.OnHTML(".area .mainInfo", func(element *colly.HTMLElement) {
		area := strings.Replace(element.Text, "平米", "", 1)
		iArea, err := strconv.Atoi(area)
		if err != nil {
			iArea = 0
		}
		db.Update(element.Request.URL.String(), bson.M{"area": iArea, "detailCrawlTime": time.Now()})
	})

	c.OnHTML("title", func(element *colly.HTMLElement) {
		logger.Sugar.Info(element.Text)
	})

	c.OnHTML(".aroundInfo .communityName .info", func(element *colly.HTMLElement) {
		db.Update(element.Request.URL.String(), bson.M{"xiaoqu": element.Text, "detailCrawlTime": time.Now()})
	})

	c.OnHTML(".l-txt", func(element *colly.HTMLElement) {
		res := strings.Replace(element.Text, "二手房", "", 99)
		res = strings.Replace(res, " ", "", 99)
		address := strings.Split(res, ">")
		db.Update(element.Request.URL.String(), bson.M{"address": address[1 : len(address)-1], "detailCrawlTime": time.Now()})
	})

	c.OnHTML(".transaction li", func(element *colly.HTMLElement) {
		if element.ChildText("span:first-child") == "挂牌时间" {
			sGTime := element.ChildText("span:last-child")
			ttime, err := time.Parse("2006-01-02", sGTime)

			if err != nil {
				ttime = time.Now()
			}

			db.Update(element.Request.URL.String(), bson.M{"zq_detail_status": 1, "guapaitime": ttime, "detailCrawlTime": time.Now()})
		}
	})

	//c.OnRequest(func(r *colly.Request) {
	//	logger.Sugar.Info("详情抓取：", r.URL.String())
	//})

	client := db.GetClient()
	ctx := db.GetCtx()

	odb := client.Database(configInfo["dbDatabase"].(string))
	lianjia := odb.Collection(configInfo["dbCollection"].(string))

	//读取出全部需要抓取详情的数据
	cur, err := lianjia.Find(ctx, bson.M{"zq_detail_status": 0})

	if err != nil {
		logger.Sugar.Error(err)
	} else {
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var item bson.M
			if err := cur.Decode(&item); err != nil {
				logger.Sugar.Errorf("数据库读取失败:%s", err.Error())
			} else {
				sucnum++
				c.Visit(item["Link"].(string))
			}

		}
	}
	return sucnum
}

func StartLianjiaErshou() {
	listFlag := make(chan int) //记录列表抓取是否完成

	go func() {
		listCrawler()
		listFlag <- 1 //列表抓取完成
	}()

	//详情抓取与列表抓取都完成了，结束主线程
	<-listFlag

	// 抓详情
	count := crawlDetail()
	if count == 0 {
		logger.Sugar.Error("抓取失败,没有数据")
	} else {
		logger.Sugar.Infof("[2/2][%d/%d] 抓取详情完成", count, count)
	}
}
