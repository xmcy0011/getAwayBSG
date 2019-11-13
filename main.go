package main

import (
	"flag"
	"github.com/getAwayBSG/configs"
	"github.com/getAwayBSG/entrance"
	"github.com/getAwayBSG/logger"
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
	flag.BoolVar(&lianjiaErshou, "lianjiaErshou", false, "抓取链家二手房数据")
	flag.BoolVar(&lianjiaZufang, "lianjiaZufang", false, "抓取链家租房数据")
	flag.BoolVar(&zhilian, "zhilian", false, "抓取智联招聘数据")
	flag.BoolVar(&clean, "clean", false, "清理缓存")
	flag.BoolVar(&info, "info", false, "保存抓取状态")
	flag.StringVar(&infoSaveTo, "info_save_to", "./numlog.txt", "输入状态文件保存位置")
}

func main() {
	logger.InitLogger("log/log.log", "debug")
	defer logger.Logger.Sync()

	flag.Parse()
	//初始化配置信息，同时输出配置信息
	if config != "" {
		configs.SetConfig(config)
	}
	logger.Sugar.Info(configs.Config())

	//进入不同入口
	if help {
		flag.Usage()
	} else if lianjiaErshou {
		entrance.Start_lianjia_ershou()
	} else if lianjiaZufang {
		entrance.Start_LianjiaZufang()
	} else if zhilian {
		entrance.Start_zhilian()
	} else if clean {
		entrance.Start_clean()
	} else if info {
		entrance.Start_info(infoSaveTo)
	} else {
		flag.Usage()
	}

}
