package configs

import (
	"github.com/getAwayBSG/logger"
	"github.com/micro/go-micro/config"
	"os"
	"path/filepath"
)

type singleton struct {
	configInfo map[string]interface{}
}

var instance *singleton
var config_path string

func GetInstance() *singleton {
	if instance == nil {
		if config_path == "" {
			instance = new(singleton)
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			err := config.LoadFile(dir + "/config.yaml")
			if err != nil {
				err = config.LoadFile("./config.yaml")
				if err != nil {
					logger.Sugar.Error("加载配置文件错误！！请确认当前目录下包含config.yaml文件或者指定配置文件参数")
					logger.Sugar.Fatal(err)
				}
			}
			conf := config.Map()
			instance.configInfo = conf
		} else {
			instance = new(singleton)
			err := config.LoadFile(config_path)
			if err != nil {
				logger.Sugar.Error("加载配置文件错误！！请确认当前目录下包含config.yaml文件或者指定配置文件参数")
				logger.Sugar.Fatal(err)
			}
			conf := config.Map()
			instance.configInfo = conf
		}

	}
	return instance
}

func Config() map[string]interface{} {

	return GetInstance().configInfo
}

func SetConfig(path string) {
	config_path = path
}
