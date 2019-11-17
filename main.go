package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/getAwayBSG/configs"
	"github.com/getAwayBSG/entrance"
	"github.com/getAwayBSG/logger"
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
	entrance.Start_clean()

	//进入不同入口
	if lianjiaErshou {
		logger.Sugar.Infof("抓取链家二手房数据,存储tableName=%s", configs.ConfigInfo.DbDatabase)
		entrance.StartLJSecondHandHouse()
	} else if lianjiaZufang {
		entrance.Start_LianjiaZufang()
	} else if zhilian {
		entrance.Start_zhilian()
	} else {
		fmt.Println("每次抓取都是全量的!")
		flag.Usage()
	}
	//else if clean {
	//	entrance.Start_clean()
	//} else if info {
	//	entrance.Start_info(infoSaveTo)
	//}
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
