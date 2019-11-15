package entrance

import (
	"context"
	"fmt"
	"github.com/getAwayBSG/configs"
	"github.com/getAwayBSG/db"
	"github.com/getAwayBSG/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
	"time"
)

func Start_clean() {
	var choice int
	if strings.Index(strings.Join(os.Args, ""), "lianjia_ershou") > -1 {
		choice = 1
	} else if strings.Index(strings.Join(os.Args, ""), "zhilian") > -1 {
		choice = 2
	} else if strings.Index(strings.Join(os.Args, ""), "lianjia_zufang") > -1 {
		choice = 3
	} else {
		logger.Sugar.Info("清除抓取状态（不清除状态的话爬虫会从上次停止位置继续抓取）")
		logger.Sugar.Info("请选择需要清除哪个爬虫的的状态数据：（输入数字）")
		logger.Sugar.Info("1.链家二手房")
		logger.Sugar.Info("2.智联")
		logger.Sugar.Info("3.链家租房")
		fmt.Scanln(&choice)

	}

	if choice == 1 {
		clean_visit()
		logger.Sugar.Info("Clear colly cookie done!")
	} else if choice == 2 {
		db.SetZhilianStatus(0, 0)
		logger.Sugar.Info("Done!")
	} else if choice == 3 {
		db.SetLianjiaZuFangStatus(0)
		clean_visit()
		logger.Sugar.Info("Done!")
	} else {
		logger.Sugar.Info("选择错误！")
	}

}

func clean_visit() {
	conf := configs.Config()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.NewClient(options.Client().ApplyURI(conf["dburl"].(string) + "/colly"))
	if err := client.Connect(ctx); err != nil {
		logger.Sugar.Error(err)
	}

	odb := client.Database("colly")
	cookies := odb.Collection("colly_cookies")
	visit := odb.Collection("colly_visited")
	//清除全部的cookies
	_, err := cookies.DeleteMany(ctx, bson.M{})
	if err != nil {
		logger.Sugar.Error(err)
	}
	_, err = visit.DeleteMany(ctx, bson.M{})
	if err != nil {
		logger.Sugar.Error(err)
	}

}
