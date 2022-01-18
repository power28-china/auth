package mongo

import (
	"context"

	"github.com/power28-china/auth/config"
	"github.com/power28-china/auth/utils/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// Client the mongodb client
	client *mongo.Client
	// DB the mongodb database.
	db  *mongo.Database
	err error
)

func init() {
	url := config.Config("MONGODB_URL")
	// logger.Sugar.Infof("loaded mongodb url from env file: %s", config.Config("MONGODB_URL"))

	if client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(url).SetMaxPoolSize(20)); err != nil {
		logger.Sugar.Fatal(err)
	}

	db = client.Database(config.Config("DATABASE_NAME"))

	// Test connection
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		logger.Sugar.Fatal(err)
	}
	logger.Sugar.Info("Connection Opened to Mongo database successfully.")
}
