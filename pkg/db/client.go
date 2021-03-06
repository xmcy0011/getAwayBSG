package db

import (
	"context"
	"getAwayBSG/pkg/configs"
	"getAwayBSG/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type singleton struct {
	client *mongo.Client
	ctx    context.Context
}

var instance *singleton

func GetInstance() *singleton {
	if instance == nil {
		instance = new(singleton)
		client, _ := mongo.NewClient(options.Client().ApplyURI(configs.ConfigInfo.DbRrl + "/" + configs.ConfigInfo.DbDatabase))
		ctx := context.Background()
		instance.client = client
		instance.ctx = ctx
		err := client.Connect(ctx)
		if err != nil {
			logger.Sugar.Fatalf("数据库连接失败！%s", err.Error())
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			logger.Sugar.Fatalf("ping error:%s", err.Error())
		}

	}
	return instance
}

func GetZhilianStatus() (int, int) {
	client := GetInstance().client
	ctx := GetInstance().ctx
	db := client.Database(configs.ConfigInfo.DbDatabase)
	lianjia_status := db.Collection("zhilian_status")
	var res bson.M

	var city_index int
	var kw_index int
	err := lianjia_status.FindOne(ctx, bson.M{}).Decode(&res)
	if err != nil {
		return 0, 0
	}
	if res["city_index"] != nil {
		city_index = int(res["city_index"].(int32))
	}

	if res["kw_index"] != nil {
		kw_index = int(res["kw_index"].(int32))
	}

	return city_index, kw_index
}

func SetZhilianStatus(cityIndex int, kwIndex int) {
	client := GetInstance().client
	ctx := GetInstance().ctx
	db := client.Database(configs.ConfigInfo.DbDatabase)
	lianjia_status := db.Collection("zhilian_status")
	lianjia_status.DeleteMany(ctx, bson.M{})
	lianjia_status.InsertOne(ctx, bson.M{"city_index": cityIndex, "kw_index": kwIndex})
}

func GetLianjiaZuFangStatus() int {
	client := GetInstance().client
	ctx := GetInstance().ctx
	db := client.Database(configs.ConfigInfo.DbDatabase)
	lianjia_status := db.Collection("lianjiazf_status")
	var res bson.M
	err := lianjia_status.FindOne(ctx, bson.M{}).Decode(&res)
	if err != nil {
		return 0
	}

	index := res["index"].(int32)
	return int(index)
}

func SetLianjiaZuFangStatus(i int) {
	client := GetInstance().client
	ctx := GetInstance().ctx
	db := client.Database(configs.ConfigInfo.DbDatabase)
	lianjia_status := db.Collection("lianjiazf_status")
	lianjia_status.DeleteMany(ctx, bson.M{})
	lianjia_status.InsertOne(ctx, bson.M{"index": i})
}

func GetCtx() context.Context {
	return GetInstance().ctx
}

func GetClient() *mongo.Client {
	return GetInstance().client
}
