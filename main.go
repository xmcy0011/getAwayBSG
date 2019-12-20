package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/getAwayBSG/entrance"
	"github.com/getAwayBSG/pkg/configs"
	"github.com/getAwayBSG/pkg/logger"
	"github.com/getwe/figlet4go"
)

// 申明配置变量
var (
	help          bool
	config        string
	lianjiaErshou bool
	lianjiaZufang bool
	zhilian       bool
	clean         bool
	info          bool
	infoSaveTo    string
)

func init() {
	flag.BoolVar(&help, "help", false, "显示帮助")
	flag.StringVar(&config, "config", "./config.yaml", "设置配置文件")
	flag.BoolVar(&lianjiaErshou, "lianjia_ershou", false, "抓取链家二手房数据")
	flag.BoolVar(&lianjiaZufang, "lianjia_zufang", false, "抓取链家租房数据")
	flag.BoolVar(&zhilian, "zhilian", false, "抓取智联招聘数据")
}

func main() {
	logger.InitLogger("log/log.log", "debug")
	defer logger.Logger.Sync()
	printIcon()

	flag.Parse()
	//初始化配置信息，同时输出配置信息
	if config != "" {
		configs.LoadConfig(config)
	} else {
		configs.LoadConfig("config.yaml")
	}
	logger.Sugar.Infof("dbAddress=%s,dbName=%s", configs.ConfigInfo.DbRrl, configs.ConfigInfo.DbDatabase)

	//进入不同入口
	if lianjiaErshou {
		entrance.CleanVisit()
		logger.Sugar.Infof("抓取链家二手房数据,存储tableName=%s", configs.ConfigInfo.DbDatabase)
		entrance.StartLJSecondHandHouse(true)
	} else if lianjiaZufang {
		entrance.Start_LianjiaZufang()
	} else if zhilian {
		entrance.Start_zhilian()
	} else {
		fmt.Println("每次抓取都是全量的!")
		flag.Usage()

		choice()
	}
	//else if clean {
	//	entrance.Start_clean()
	//} else if info {
	//	entrance.Start_info(infoSaveTo)
	//}
}

func choice() {
	var choice int
	for ; ; {
		logger.Sugar.Info("请选择任务")
		logger.Sugar.Info("1.爬取链家二手房（全量）")
		logger.Sugar.Info("2.爬取链家二手房（详情）")
		logger.Sugar.Info("3.链家二手房补充经纬度")
		//logger.Sugar.Info("3.爬智联")
		//logger.Sugar.Info("4.链家租房")

		_, err := fmt.Scanln(&choice)
		if err == nil {
			if choice == 1 {
				entrance.CleanVisit()
				logger.Sugar.Infof("抓取链家二手房数据,存储tableName=%s", configs.ConfigInfo.DbDatabase)
				entrance.StartLJSecondHandHouse(true)
				break
			} else if choice == 2 {
				entrance.CleanVisit()
				logger.Sugar.Infof("抓取链家二手房数据,存储tableName=%s", configs.ConfigInfo.DbDatabase)
				entrance.StartLJSecondHandHouse(false)
				break
			} else if choice == 3 {
				logger.Sugar.Infof("链家二手房数据补充经纬度信息,房源tableName=%s", configs.ConfigInfo.DbDatabase)
				entrance.StartGeocodeLJ()
				break
			} else {
				logger.Sugar.Info("选择错误！")
			}
		} else {
			logger.Sugar.Error("选择错误:" + err.Error())
		}
	}
}

func printIcon() {
	appName := "getAwayBSG"
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render(appName)
	// 黑白
	//fmt.Println(renderStr)

	colors := [...]color.Attribute{
		color.FgMagenta,
		color.FgYellow,
		color.FgBlue,
		color.FgCyan,
		color.FgRed,
		color.FgWhite,
	}
	options := figlet4go.NewRenderOptions()
	options.FontColor = make([]color.Attribute, len(appName))
	for i := range options.FontColor {
		options.FontColor[i] = colors[i%len(colors)]
	}
	renderStr, _ = ascii.RenderOpts(appName, options)
	// 彩色
	fmt.Println(renderStr)
}
