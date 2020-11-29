package db

import (
	"context"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	_id        primitive.ObjectID
	Username   string
	Password   string // hashed
	Created    int64  // unix
	Jwt        string
	LoggedOn   bool
	LastOnline int64 // unix
}

func (user *User) CreateUser(context context.Context) (err error) {
	_, err = db.Collection("users").InsertOne(context, user)
	if err != nil {
		glog.Error(err)
		return
	}
	return
}
