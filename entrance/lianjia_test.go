package entrance

import (
	"getAwayBSG/pkg/logger"
	"testing"
)

func TestGetSecondAreas(t *testing.T) {
	logger.InitLogger("log/log.log", "debug")
	defer logger.Logger.Sync()

	arr, count, err := getSecondAreas("https://sh.lianjia.com/ershoufang/pudong/")
	if err == nil {
		logger.Sugar.Infof("浦东中房源套数:%d", count)

		for i := range arr {
			logger.Sugar.Infof("name:%s,url:%s", arr[i].Name, arr[i].Url)
		}
	}
}
