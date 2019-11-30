package configs

import (
	"github.com/getAwayBSG/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 租房
type RentCityConfig struct {
	Link string `yaml:"link"`
	Name string `yaml:"name"`
}

// 智联招聘
type ZLJob struct {
	Name     string `yaml"name"`
	Url      string `yaml:"url"`
	Code     int    `yaml:"code"`
	Pinyin   string `yaml:"pinyin"`
	Priority int    `yaml:"priority"`
}

type CrawlerConfig struct {
	DbRrl                  string `yaml:"dbUrl"`
	DbDatabase             string `yaml:"dbDatabase"`
	CollyDatabase          string `yaml:"collyDatabase"`
	RentCollection         string `yaml:"rentCollection"`
	CrawlDelay             int    `yaml:"crawlDelay"`
	CrawlDetailRoutineNum  int    `yaml:"crawlDetailRoutineNum"`
	TwoHandHouseCollection string `yaml:"twoHandHouseCollection"`
	ZlDBCollection         string `yaml:"zlDBCollection"`

	ProxyList []string `yaml:"proxyList"`

	ZlKeyWords   []string         `yaml:"zlKeyWords"`
	CityList     []string         `yaml:"cityList"`
	RentCityList []RentCityConfig `yaml:"rentCityList"`
	ZlCityList   []ZLJob          `yaml:"zlCityList"`
}

// 配置文件信息
var ConfigInfo = &CrawlerConfig{}

// 加载配置文件
func LoadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Sugar.Error(err.Error())
	}

	err = yaml.Unmarshal(data, ConfigInfo)
	if err != nil {
		logger.Sugar.Error(err.Error())
	}
}
