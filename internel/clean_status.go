package internel

import (
	"context"
	"getAwayBSG/pkg/configs"
	"getAwayBSG/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func CleanVisit() {
	logger.Sugar.Info("clear colly cookies and visited...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.NewClient(options.Client().ApplyURI(configs.ConfigInfo.DbRrl + "/colly"))
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
