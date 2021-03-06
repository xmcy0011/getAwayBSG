package old

import (
	"getAwayBSG/pkg/configs"
	"getAwayBSG/pkg/db"
	"getAwayBSG/pkg/logger"
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

func TcrawlerOneCityZuFang(cityUrl string, cityname string) {
	c := colly.NewCollector()

	if configs.ConfigInfo.CrawlDelay > 0 {
		_ = c.Limit(&colly.LimitRule{
			DomainGlob: "*",
			Delay:      time.Duration(configs.ConfigInfo.CrawlDelay) * time.Second,
		})
	}

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
	c.OnRequest(func(r *colly.Request) {
		logger.Sugar.Info("列表抓取：", r.URL.String())
	})

	c.OnHTML("title", func(element *colly.HTMLElement) {
		logger.Sugar.Info(element.Text)
	})

	c.OnHTML(".content__list--item", func(element *colly.HTMLElement) {
		var err error
		var link string
		var title string
		var address string
		var area string
		var price int
		var mianji int
		element.ForEach(".twoline a", func(i int, element *colly.HTMLElement) {
			link = "https://" + element.Request.URL.Host + element.Attr("href")
			title = strings.TrimSpace(element.Text)
		})

		element.ForEach(".content__list--item--des a", func(i int, element *colly.HTMLElement) {
			if i == 0 {
				address = element.Text
			} else {
				area = element.Text
			}
		})

		desc := element.ChildText(".content__list--item--des")
		desc = strings.ReplaceAll(desc, " ", "")
		desc = strings.ReplaceAll(desc, "\n", "")
		logger.Sugar.Info(desc)
		re, _ := regexp.Compile("(\\d+)㎡/")
		indexs := re.FindStringIndex(desc)
		if len(indexs) == 2 {

			mianji, err = strconv.Atoi(desc[indexs[0] : indexs[1]-4])
			if err != nil {
				mianji = 0
			}
		} else {
			mianji = 0
		}

		element.ForEach(".content__list--item-price em", func(i int, element *colly.HTMLElement) {
			var err error
			price, err = strconv.Atoi(element.Text)
			if err != nil {
				price = 0
			}
		})

		//fmt.Println(price)
		//fmt.Println(link)
		//fmt.Println(title)
		//fmt.Println(address)
		//fmt.Println(area)
		//fmt.Println(cityname)
		logger.Sugar.Info("--------------------")

		client := db.GetClient()
		ctx := db.GetCtx()

		db := client.Database(configs.ConfigInfo.DbDatabase)
		lianjia := db.Collection(configs.ConfigInfo.RentCollection)
		_, err = lianjia.InsertOne(ctx, bson.M{
			"Link":       link,
			"title":      title,
			"address":    address,
			"area":       area,
			"price":      price,
			"city":       cityname,
			"mianji":     mianji,
			"crawl_time": time.Now(),
		})
		if err != nil {
			if !strings.Contains(err.Error(), "multiple write errors") {
				logger.Sugar.Error("数据库插入失败:%s", err.Error())
			}
		}

	})

	c.OnHTML(".content__pg", func(element *colly.HTMLElement) {
		totalPage := element.Attr("data-totalpage")
		iTotalPage, err := strconv.Atoi(totalPage)

		if err == nil {
			for i := 2; i < iTotalPage; i++ {

				re, _ := regexp.Compile("pg\\d+/*")
				goUrl := re.ReplaceAllString(element.Request.URL.String(), "")

				err = c.Visit(goUrl + "pg" + strconv.Itoa(i) + "/")
				if err != nil && err.Error() != "URL already visited" {
					logger.Sugar.Error(err)
				}

			}
		}
	})

	c.OnHTML(".filter ul", func(element *colly.HTMLElement) {

		data_target := element.Attr("data-target")

		if data_target == "area" {
			element.ForEach("a", func(i int, element *colly.HTMLElement) {
				//// 切换地点
				u, err := url.Parse(cityUrl)
				if err != nil {
					panic(err)
				}
				rootUrl := u.Scheme + "://" + u.Host
				goUrl := element.Attr("href")
				u, err = url.Parse(goUrl)
				if err != nil && err.Error() != "URL already visited" {
					logger.Sugar.Error(err)
				}
				if u.Scheme == "" {
					goUrl = rootUrl + u.Path
				} else {
					goUrl = u.String()
				}
				re, _ := regexp.Compile("pg\\d+/*")
				goUrl = re.ReplaceAllString(goUrl, "")
				err = c.Visit(goUrl)
				if err != nil && err.Error() != "URL already visited" {
					logger.Sugar.Error(err)
				}

			})
		}

	})

	err := c.Visit(cityUrl)
	if err != nil && err.Error() != "URL already visited" {
		logger.Sugar.Error(err)
	}

}

func Start_LianjiaZufang() {
	cityList := configs.ConfigInfo.RentCityList

	for i := db.GetLianjiaZuFangStatus(); i < len(cityList); i++ {
		TcrawlerOneCityZuFang(cityList[i].Link, cityList[i].Name)
		db.SetLianjiaZuFangStatus(i)
	}
}
