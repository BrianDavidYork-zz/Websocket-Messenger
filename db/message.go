package db

import (
	"context"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	_id            primitive.ObjectID
	ConversationId primitive.ObjectID
	Message        string
	State          int // 0 = normal; 1 = edited; 2 = deleted
	Sender         string
	Created        int64 // unix
}

func (msg *Message) Create(context context.Context) (mid primitive.ObjectID, err error) {
	msg.Created = time.Now().Unix()
	result, err := db.Collection("messages").InsertOne(context, msg)
	if err != nil {
		glog.Error(err)
		return
	}
	mid = result.InsertedID.(primitive.ObjectID)
	return
}

func GetMessageById(context context.Context, messageId primitive.ObjectID) (msg Message, err error) {
	err = db.Collection("messages").FindOne(context, bson.M{"_id": messageId}).Decode(&msg)
	if err != nil {
		glog.Error(err)
		return
	}
	return
}

func EditMessage(context context.Context, messageId primitive.ObjectID, msg string) (err error) {
	_, err = db.Collection("messages").UpdateOne(context, bson.M{"_id": messageId}, bson.M{"$set": bson.M{"message": msg, "state": 1}})
	if err != nil {
		glog.Error(err)
		return
	}
	return
}
