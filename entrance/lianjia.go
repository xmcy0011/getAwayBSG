package entrance

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/getAwayBSG/pkg/configs"
	"github.com/getAwayBSG/pkg/db"
	"github.com/getAwayBSG/pkg/logger"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	cachemongo "github.com/zolamk/colly-mongo-storage/colly/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	AreaUrl   string
}

type HouseInfo struct {
	Title      string `json:"title"`
	Link       string `json:"link"`
	TotalPrice int32  `json:"total_price"`
}

// 获取一个城市下所有区域的url
func getAreaUrl(cityUrl string) []string {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	areaArr := make([]string, 0)

	c.OnHTML("body", func(element *colly.HTMLElement) {
		element.ForEachWithBreak(".position div div", func(i int, element *colly.HTMLElement) bool {
			u, err := url.Parse(cityUrl)
			if err != nil {
				panic(err)
			}

			element.ForEach("a", func(i int, element *colly.HTMLElement) {
				goUrl := element.Attr("href")
				rootUrl := u.Scheme + "://" + u.Host
				areaArr = append(areaArr, rootUrl+goUrl)
			})

			logger.Sugar.Infof("%s,抓取地区共:%d", cityUrl, len(areaArr))
			return false
		})
	})

	err := c.Visit(cityUrl)
	if err != nil {
		logger.Sugar.Fatalf("%s:%s", err.Error(), cityUrl)
	}

	return areaArr
}

func getListProgress(cityIndex int, cityCount int, areaIndexMin int, areaIndexMax, totalArea int,
	curPage int, totalPage int, totalCount int) string {
	// [1/2]
	// [%d/%d]：城市
	// [%d/%d]：分页
	return fmt.Sprintf("[1/2][%d/%d 城][%d-%d/%d 区][%d/%d 页][%d 条]", cityIndex+1, cityCount,
		areaIndexMin, areaIndexMax, totalArea, curPage, totalPage, totalCount)
}
func getDetailProgress(curCount int, totalCount int) string {
	percent := int(crawlerDetailSuccessCount) * 100 / crawlerDetailCount
	percentStr := strconv.Itoa(percent) + "%"
	return fmt.Sprintf("[2/2][%s%s][%d/%d 成功/总数]", percentStr, "%", crawlerDetailSuccessCount, crawlerDetailCount)
}

func decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func crawlerOneArea(areaUrl string, cityName string, progressChan chan *Page) {
	var areaName string
	var page = &Page{CurPage: 1, TotalPage: 1, AreaUrl: areaUrl}
	var nextPageUrl string

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

	// 绑定
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
				page.CurPage = tempPage.CurPage
				page.TotalPage = tempPage.TotalPage

				//progressInfo := getListProgress(cityIndex, cityCount, areaIndex, areaCount, page.CurPage, page.TotalPage)
				//logger.Sugar.Infof("%s[%s] totalCount=%s,totalPage=%d,curPage=%d",
				//	progressInfo, cityName, areaName, page.TotalPage, page.CurPage)
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

			// 位置
			// 中城丽景香山-武广新城
			var listVillageName = ""
			var listAreaName = ""

			positionInfo := e.ChildText(".positionInfo")
			positionInfo = strings.Replace(positionInfo, " ", "", -1)
			if temp := strings.Split(positionInfo, "-"); len(temp) >= 2 {
				listAreaName = temp[0]
				listVillageName = temp[1]
			}

			// 房屋信息 户型、大小、朝向、装修、楼层、年代（可为空）、板楼
			// 3室2厅|133.55平米|南|精装|中楼层(共33层)|2012年建|板塔结合
			var listHouseType = ""
			var listHouseSize = 0.0
			var listHouseOrientations = ""
			var listHouseDecorate = ""
			var listHouseFloor = ""
			var listHouseBorn = ""
			var listHouseWhat = ""

			houseInfo := e.ChildText(".houseInfo")
			houseInfo = strings.Replace(houseInfo, " ", "", -1)
			temp := strings.Split(houseInfo, "|")
			if len(temp) >= 6 {
				listHouseType = temp[0]
				tempSize := temp[1]
				listHouseOrientations = temp[2]
				listHouseDecorate = temp[3]
				listHouseFloor = temp[4]

				tempSize = strings.Replace(tempSize, "平米", "", 1)
				listHouseSize, err = strconv.ParseFloat(tempSize, 10)
				if err != nil {
					listHouseSize = 0
				}

				if len(temp) >= 7 {
					listHouseBorn = temp[5]
					listHouseWhat = temp[6]
				} else {
					listHouseWhat = temp[5]
				}
			}

			// 关注信息 关注人数、发布时间
			// 9人关注/2个月以前发布
			followInfo := e.ChildText(".followInfo")
			followInfo = strings.Replace(followInfo, " ", "", -1)
			// tag
			tag := make([]string, 0)
			e.ForEach(".tag", func(i int, element *colly.HTMLElement) {
				element.ForEach("span", func(i int, element *colly.HTMLElement) {
					tag = append(tag, element.Text)
				})
			})

			//progressInfo := getListProgress(cityIndex, cityCount, areaIndex, areaCount, page.CurPage, page.TotalPage)
			//logger.Sugar.Infof("%s[%d] %s,%s,%s,总价：%d 万元，每平米：%d",
			//	progressInfo, curCount, cityName, areaName, title, iPrice, iUnitPrice)

			db.Add(bson.M{
				"DetailStatus": 0, "Title": title, "TotalPrice": iPrice, "UnitPrice": iUnitPrice,
				"Link": link, "ListCrawlTime": time.Now().Format("2006-01-02 15:04:05"), "City": cityName,
				"ListVillageName": listVillageName, "ListAreaName": listAreaName, "ListHouseType": listHouseType,
				"ListHouseSize": listHouseSize, "ListHouseOrientations": listHouseOrientations, "ListHouseDecorate": listHouseDecorate,
				"ListHouseFloor": listHouseFloor, "ListHouseBorn": listHouseBorn, "ListHouseWhat": listHouseWhat, "Tag": tag,
			}, link)
		})

		// 下一页
		element.ForEach(".page-box .house-lst-page-box", func(i int, element *colly.HTMLElement) {
			var tempPage Page
			err := json.Unmarshal([]byte(element.Attr("page-data")), &tempPage)
			if err == nil {
				page.CurPage = tempPage.CurPage
				if page.CurPage < page.TotalPage {
					re, _ := regexp.Compile("pg\\d+/*")
					url := re.ReplaceAllString(element.Request.URL.String(), "")

					nextPageUrl = url + "pg" + strconv.Itoa(tempPage.CurPage+1)
				}
			}
		})
	})

	// 初始化访问
	err := c.Visit(areaUrl)
	if err != nil {
		logger.Sugar.Debugf("%s:%s", err.Error(), areaUrl)
	}
	progressChan <- page // 更新进度

	// 一个地区下的所有分页数据
	for j := page.CurPage; j < page.TotalPage; j++ {
		err = c.Visit(nextPageUrl)
		if err == nil {
			if page.CurPage != (j + 1) {
				logger.Sugar.Errorf("修正分页数据，page.CurPage=%d,j=%d，忽略继续", page.CurPage, j)
				j = page.CurPage
			}
		} else {
			logger.Sugar.Error(err)
		}
		if page.CurPage > 10 {
			break
		}
		progressChan <- page // 更新进度
	}
}

func crawlerOneCity(cityName string, cityUrl string, cityIndex int, cityCount int) {
	areaArr := getAreaUrl(cityUrl)

	maxRoutine := configs.ConfigInfo.CrawlDetailRoutineNum
	if maxRoutine > 4 {
		maxRoutine = 4
		logger.Sugar.Infof("为防止河蟹，列表抓取已修正为 %d 协程同时抓取", maxRoutine)
	}

	const pageSize = 30 // 链接，分页30条

	areaIndexMin := 0
	areaIndexMax := 0
	areaTotal := len(areaArr)
	waitGroup := &sync.WaitGroup{}
	progressChan := make(chan *Page, 0)

	// 计算4个routine同时爬所有地区，要进行几轮
	loopCount := len(areaArr) / maxRoutine
	if temp := len(areaArr) % maxRoutine; temp != 0 {
		loopCount++
	}

	// print list progress
	isComplete := false
	curPage := 0
	totalPage := 0
	totalCount := 0
	go func() {
		isInit := false
		initPageAreaMap := make(map[string]*Page, maxRoutine)

		for {
			// wait chan
			page := <-progressChan
			curPage++
			totalCount += pageSize

			// 统计所有条数，等待所有协程初始化完成后，才打印进度
			if value, ok := initPageAreaMap[page.AreaUrl]; ok {
				value.CurPage = page.CurPage
			} else {
				initPageAreaMap[page.AreaUrl] = page
				totalPage += page.TotalPage
			}
			if len(initPageAreaMap) >= maxRoutine {
				isInit = true
			}

			// 打印进度
			if isInit {
				logger.Sugar.Info(getListProgress(cityIndex, cityCount, areaIndexMin, areaIndexMax, areaTotal, curPage, totalPage, totalCount))
			}

			if isComplete {
				break
			}
		}
	}()

	// 多线程同时爬取多个地区下所有列表
	for i := 0; i < loopCount; i++ {
		curPage = 0
		totalPage = 0
		areaIndexMin = (i * maxRoutine) + 1
		if areaIndexMax = (i + 1) * maxRoutine; areaIndexMax > len(areaArr) {
			areaIndexMax = len(areaArr)
		}

		waitGroup.Add(maxRoutine)
		for routine := 0; routine < maxRoutine; routine++ {
			areaIndex := routine + i*maxRoutine
			if areaIndex < len(areaArr) {
				go func(areaUrl string, cityName string, progressChan chan *Page) {
					crawlerOneArea(areaUrl, cityName, progressChan)
					// complete one routine,notify it
					waitGroup.Done()
				}(areaArr[areaIndex], cityName, progressChan)
			} else {
				waitGroup.Done()
			}
		}

		// wait all routine done
		waitGroup.Wait()
	}

	isComplete = true
	time.Sleep(time.Second * 2)
}

func listCrawler() {
	count := len(configs.ConfigInfo.CityList)
	for i := 0; i < count; i++ {
		url := configs.ConfigInfo.CityList[i]
		name := strings.Split(url[8:], ".")[0] // https://cs.lianjia... -> cs

		logger.Sugar.Infof("[1/2][%d/%d] 抓取城市：%s,url=%s", i+1, count, name, url)
		crawlerOneCity(name, url, i, count)
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

	var baseAttr []string        // 基本属性
	var transactionAttr []string // 交易属性

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
			baseAttr = append(baseAttr, label+":"+element.Text[(index+len(label)):])
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
						if liText == "抵押信息:" {
							bettyString := strings.TrimSpace(element.Text)
							bettyString = strings.ReplaceAll(bettyString, "\\n", "")
							liText += bettyString
						} else {
							liText += element.Text
						}
					}
				})
				transactionAttr = append(transactionAttr, liText)
			}
		})
	})

	// 标题
	c.OnHTML("title", func(element *colly.HTMLElement) {
		title = element.Text
	})

	for i := range houseArr {
		baseAttr = make([]string, 0)
		transactionAttr = make([]string, 0)
		url := houseArr[i].Link
		err := c.Visit(url)
		if err != nil {
			logger.Sugar.Errorf("%s[协程%d],抓取失败:%s,url=%s", getDetailProgress(startNum+1, total),
				routineIndex, err.Error(), url)
			db.Update(url, bson.M{"DetailStatus": 2})
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
				"BeOnlineTime":    beOnlineTime.Format("2006-01-02 15:04:05"),
				"DetailCrawlTime": time.Now().Format("2006-01-02 15:04:05")})
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

func pingMongoDb() error {
	logger.Sugar.Info("ping mongoDb,timeout 10s ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.NewClient(options.Client().ApplyURI(configs.ConfigInfo.DbRrl + "/" + configs.ConfigInfo.DbDatabase))
	defer client.Disconnect(ctx)
	if err := client.Connect(ctx); err != nil {
		return err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	logger.Sugar.Info("ping mongoDb success")
	return nil
}

// 开始抓取链家二手房（先抓取概要列表、完成后再批量抓取二手房详情）
// @param crawlerList:是否抓取列表
func StartLJSecondHandHouse(crawlerList bool) {
	if err := pingMongoDb(); err != nil {
		logger.Sugar.Errorf("mongoDb connected error:%s", err.Error())
		return
	}

	if crawlerList {
		listFlag := make(chan int)
		go func() {
			logger.Sugar.Info("[1/2] 开始抓取城市二手房概要信息")
			listCrawler()
			listFlag <- 1 //列表抓取完成
		}()

		//阻塞主线程，等待列表抓取
		<-listFlag
	}

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