package entrance

import (
	"encoding/json"
	"fmt"
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
	"sync"
	"sync/atomic"
	"time"
)

var crawlerDetailCount = 0               // 总爬取详情数
var crawlerDetailSuccessCount = int32(0) // 总爬取详情成功数

type Page struct {
	TotalPage int
	CurPage   int
}

type HouseInfo struct {
	Title      string `json:"title"`
	Link       string `json:"link"`
	TotalPrice int32  `json:"total_price"`
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

func getListProgress(cityIndex int, cityCount int, areaIndex int, areaCount int, pageNum int, pageCount int) string {
	// [1/2]
	// [%d/%d]：城市
	// [%d/%d]：地区
	// [%d/%d]：分页
	return fmt.Sprintf("[1/2][%d/%d][%d/%d][%d/%d]",
		cityIndex+1, cityCount, areaIndex+1, areaCount, pageNum, pageCount)
}

func getDetailProgress(curCount int, totalCount int) string {
	percent := int(crawlerDetailSuccessCount) * 100 / crawlerDetailCount
	percentStr := strconv.Itoa(percent) + "%"
	return fmt.Sprintf("[2/2][%s][%d/%d]", percentStr, crawlerDetailSuccessCount, crawlerDetailCount)
}

func decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func crawlerOneCity(cityUrl string, cityIndex int, cityCount int) {
	c := colly.NewCollector()

	if configs.ConfigInfo.CrawlDelay > 0 {
		// randomDelay:200-delay
		err := c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Delay:       time.Millisecond * 200,
			RandomDelay: time.Duration(configs.ConfigInfo.CrawlDelay-200) * time.Millisecond,
		})
		if err != nil {
			logger.Sugar.Error(err)
		}
	}

	// 代理
	if configs.ConfigInfo.ProxyList != nil {
		rp, err := proxy.RoundRobinProxySwitcher(configs.ConfigInfo.ProxyList...)
		if err != nil {
			logger.Sugar.Error(err)
		}
		c.SetProxyFunc(rp)
	}

	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	storage := &cachemongo.Storage{
		Database: "colly",
		URI:      configs.ConfigInfo.DbRrl + "/colly",
	}
	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}

	areaListChan := make(chan []string, 1)
	getAeraUrl(cityUrl, areaListChan)
	areaArr := <-areaListChan

	var page Page
	var nextPageUrl string
	var areaName string
	var areaCount = len(areaArr)

	arr := strings.Split(cityUrl, ".")
	cityName := "unknown_city"
	if len(arr) > 1 {
		cityName = arr[0][8:]
	}

	// 挨个爬各个地区的房源
	for areaIndex := range areaArr {
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

					progressInfo := getListProgress(cityIndex, cityCount, areaIndex, areaCount, page.CurPage, page.TotalPage)
					logger.Sugar.Infof("%s[%s] totalCount=%s,totalPage=%d,curPage=%d",
						progressInfo, cityName, areaName, page.TotalPage, page.CurPage)
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

				progressInfo := getListProgress(cityIndex, cityCount, areaIndex, areaCount, page.CurPage, page.TotalPage)
				logger.Sugar.Infof("%s[%d] %s,%s,%s,总价：%d 万元，每平米：%d",
					progressInfo, curCount, cityName, areaName, title, iPrice, iUnitPrice)

				db.Add(bson.M{"DetailStatus": 0, "Title": title, "TotalPrice": iPrice, "UnitPrice": iUnitPrice, "Link": link, "ListCrawlTime": time.Now()})
			})

			// 下一页
			element.ForEach(".page-box .house-lst-page-box", func(i int, element *colly.HTMLElement) {
				var tempPage Page
				err := json.Unmarshal([]byte(element.Attr("page-data")), &tempPage)
				if err == nil {
					page = tempPage
					if page.CurPage < page.TotalPage {
						re, _ := regexp.Compile("pg\\d+/*")
						url := re.ReplaceAllString(element.Request.URL.String(), "")

						nextPageUrl = url + "pg" + strconv.Itoa(tempPage.CurPage+1)
					}
				}
			})
		})

		page.CurPage = 1
		page.TotalPage = 1
		err := c.Visit(areaArr[areaIndex])
		if err != nil {
			logger.Sugar.Debugf("%s:%s", err.Error(), cityUrl)
		}

		// 一个地区下的所有分页数据
		for j := page.CurPage; j < page.TotalPage; j++ {
			err = c.Visit(nextPageUrl)
			if err != nil {
				logger.Sugar.Info(err)
			}

			if page.CurPage != (j + 1) {
				logger.Sugar.Errorf("修正分页数据，page.CurPage=%d,j=%d，忽略继续", page.CurPage, j)
				j = page.CurPage
			}
		}
	}
}

func listCrawler() {
	count := len(configs.ConfigInfo.CityList)
	for i := 0; i < count; i++ {
		cityName := configs.ConfigInfo.CityList[i]
		logger.Sugar.Infof("[1/2][%d/%d] 抓取城市：%s", i+1, count, cityName)
		crawlerOneCity(cityName, i, count)
	}
}

func crawlerOneDetail(startNum int, routineIndex int, houseArr []HouseInfo, total int) {
	c := colly.NewCollector()

	//设置延时
	if configs.ConfigInfo.CrawlDelay > 0 {
		err := c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Delay:       200 * time.Millisecond,
			RandomDelay: time.Millisecond * time.Duration(configs.ConfigInfo.CrawlDelay-200),
		})
		if err != nil {
			logger.Sugar.Error(err)
		}
	}

	//设置代理
	if configs.ConfigInfo.ProxyList != nil {
		rp, err := proxy.RoundRobinProxySwitcher(configs.ConfigInfo.ProxyList...)
		if err != nil {
			logger.Sugar.Error(err)
		}
		c.SetProxyFunc(rp)
	}

	//随机UA
	extensions.RandomUserAgent(c)
	//自动referer
	extensions.Referer(c)
	//设置MongoDB存储状态信息
	storage := &cachemongo.Storage{
		Database: "colly",
		URI:      configs.ConfigInfo.DbRrl + "/colly",
	}
	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}

	var roomInfo string  // 户型,3室1厅
	var floorInfo string // 楼层,低楼层/共17层

	var directionInfo string // 朝向,南北
	var decorateInfo string  // 装修,平层/精装

	var size float64         // 大小,平米
	var completedInfo string // 竣工时间,2007年建/板楼

	var villageName string   // 小区名称
	var areaName []string    // 所在区域
	var houseRecordLJ string // 链家房源编号

	var title string // 标题

	var baseAttr string        // 基本属性
	var transactionAttr string // 交易属性

	var beOnlineTime time.Time // 挂牌时间

	// 户型+楼层
	c.OnHTML(".houseInfo", func(element *colly.HTMLElement) {
		element.ForEach(".mainInfo", func(i int, element *colly.HTMLElement) {
			roomInfo = element.Text
		})
		element.ForEach(".subInfo", func(i int, element *colly.HTMLElement) {
			floorInfo = element.Text
		})
	})

	// 朝向+装修
	c.OnHTML(".type", func(element *colly.HTMLElement) {
		element.ForEach(".mainInfo", func(i int, element *colly.HTMLElement) {
			directionInfo = element.Text
		})
		element.ForEach(".subInfo", func(i int, element *colly.HTMLElement) {
			decorateInfo = element.Text
		})
	})

	// 大小+竣工时间
	c.OnHTML(".area", func(element *colly.HTMLElement) {
		element.ForEach(".mainInfo", func(i int, element *colly.HTMLElement) {
			area := strings.Replace(element.Text, "平米", "", 1)
			value, err := strconv.ParseFloat(area, 32)
			if err != nil {
				value = 0
			}
			size = decimal(value) // 保留2位小数
		})
		element.ForEach(".subInfo", func(i int, element *colly.HTMLElement) {
			completedInfo = element.Text
		})
	})

	// 小区名称
	c.OnHTML(".aroundInfo .communityName .info", func(element *colly.HTMLElement) {
		villageName = element.Text
	})

	// 所在区域
	c.OnHTML(".l-txt", func(element *colly.HTMLElement) {
		res := strings.Replace(element.Text, "二手房", "", 99)
		res = strings.Replace(res, " ", "", 99)
		areaName = strings.Split(res, ">")
	})

	// 房源编号
	c.OnHTML(".houseRecord .info", func(element *colly.HTMLElement) {
		arr := strings.Split(element.Text, "举")
		houseRecordLJ = arr[0]
	})

	// 基本属性
	c.OnHTML(".base .content", func(element *colly.HTMLElement) {
		element.ForEach("li", func(i int, element *colly.HTMLElement) {
			var label = ""
			element.ForEach("span", func(i int, element *colly.HTMLElement) {
				label = element.Text
			})
			index := strings.Index(element.Text, label)
			baseAttr += label + ":" + element.Text[(index+len(label)):] + "|"
		})
	})

	// 交易属性
	c.OnHTML(".transaction .content", func(element *colly.HTMLElement) {
		element.ForEach("li", func(i int, element *colly.HTMLElement) {
			// 挂牌时间
			if element.ChildText("span:first-child") == "挂牌时间" {
				sGTime := element.ChildText("span:last-child")
				parsedTime, err := time.Parse("2006-01-02", sGTime)
				if err != nil {
					parsedTime = time.Now()
				}
				beOnlineTime = parsedTime
			} else {
				var liText = ""
				element.ForEach("span", func(i int, element *colly.HTMLElement) {
					if i == 0 {
						liText = element.Text + ":"
					} else {
						liText += element.Text
					}
				})
				transactionAttr += liText + "|"
			}
		})
	})

	// 标题
	c.OnHTML("title", func(element *colly.HTMLElement) {
		title = element.Text
	})

	for i := range houseArr {
		url := houseArr[i].Link
		err := c.Visit(url)
		if err != nil {
			logger.Sugar.Errorf("%s[协程%d],抓取失败:%s,url=%s", getDetailProgress(startNum+1, total),
				routineIndex, err.Error(), url)
			db.Update(url, bson.M{"DetailStatus": 1})
		} else {
			// 原子操作，多线程安全
			atomic.AddInt32(&crawlerDetailSuccessCount, 1)
			logger.Sugar.Infof("%s[协程%d],标题:%s,价格:%d,房源编号:%s,朝向:%s,装修:%s", getDetailProgress(startNum+1, total),
				routineIndex, title, houseArr[i].TotalPrice, houseRecordLJ, directionInfo, decorateInfo)

			db.Update(url, bson.M{
				"DetailStatus":    1,
				"RoomInfo":        roomInfo,
				"FloorInfo":       floorInfo,
				"DirectionInfo":   directionInfo,
				"DecorateInfo":    decorateInfo,
				"Size":            size,
				"CompletedInfo":   completedInfo,
				"VillageName":     villageName,
				"AreaName":        areaName,
				"HouseRecordLJ":   houseRecordLJ,
				"BaseAttr":        baseAttr,
				"TransactionAttr": transactionAttr,
				"BeOnlineTime":    beOnlineTime,
				"DetailCrawlTime": time.Now()})
		}
		startNum++
	}
}

func crawlerDetail() {
	var routineCount int = 0

	client := db.GetClient()
	ctx := db.GetCtx()

	odb := client.Database(configs.ConfigInfo.DbDatabase)
	dbCollection := odb.Collection(configs.ConfigInfo.TwoHandHouseCollection)

	//读取出全部需要抓取详情的数据
	cur, err := dbCollection.Find(ctx, bson.M{"DetailStatus": 0})
	if err != nil {
		logger.Sugar.Fatalf("数据库读取失败:", err.Error())
		return
	}
	var houseArr = make([]HouseInfo, 0)
	err = cur.All(ctx, &houseArr)
	if err != nil {
		logger.Sugar.Fatalf("数据库读取失败:", err.Error())
		return
	}
	crawlerDetailCount = len(houseArr)
	defer cur.Close(ctx)

	routineCount = configs.ConfigInfo.CrawlDetailRoutineNum
	logger.Sugar.Infof("[2/2] 开始抓取二手房详情,总数=%d,并行抓取协程数=%d", crawlerDetailCount, routineCount)

	var wg sync.WaitGroup
	for j := 0; j < int(routineCount); j++ {
		perCount := crawlerDetailCount / routineCount
		var tempHouseArr []HouseInfo
		var startCount = j * perCount
		var endCount int
		if (j + 1) == int(routineCount) {
			endCount = crawlerDetailCount
			tempHouseArr = houseArr[startCount:crawlerDetailCount] // 除不尽的，全部交给最后一个协程
		} else {
			endCount = (j + 1) * perCount
			tempHouseArr = houseArr[startCount:endCount]
		}

		wg.Add(1)
		go func(startNum int, routineIndex int, houseArr []HouseInfo) {
			defer wg.Add(-1)
			// 1协程抓取一组数据
			crawlerOneDetail(startNum, routineIndex, tempHouseArr, crawlerDetailCount)
		}(j*perCount, j, tempHouseArr)
		logger.Sugar.Infof("[2/2] 第 %d 组协程抓取 [%d-%d] 的房屋详情", j+1, startCount, endCount)
	}

	wg.Wait() // 等待所有协程完成
}

func StartLJSecondHandHouse() {
	listFlag := make(chan int)
	go func() {
		logger.Sugar.Info("[1/2] 开始抓取城市二手房概要信息")
		listCrawler()
		listFlag <- 1 //列表抓取完成
	}()

	//阻塞主线程，等待列表抓取
	<-listFlag

	// 抓详情
	logger.Sugar.Info("[2/2] 开始抓取城市二手房详细信息")
	time.Sleep(time.Second * 3)
	crawlerDetail()
	if crawlerDetailSuccessCount == 0 {
		logger.Sugar.Error("[2/2] 抓取详情失败,没有数据，结束二手房抓取!")
	} else {
		logger.Sugar.Infof("[2/2] 抓取详情完成，成功数=%d,总数=%d，结束二手房抓取!", crawlerDetailSuccessCount, crawlerDetailCount)
	}
}
