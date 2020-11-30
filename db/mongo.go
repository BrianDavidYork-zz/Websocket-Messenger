package db

import (
	"context"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// package-scoped declaration
var db *mongo.Database

func StartMongo() {
	// The code in this function comes from the go mongo-driver docs
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	db = client.Database("Messenger")

	// unique username index
	_, err = db.Collection("users").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		glog.Info(err)
		panic(err)
	}

	glog.Info("MongoDb started")
}
