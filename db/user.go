package db

import (
	"context"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
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

type Login struct {
	Username string
	Password string
}

func (user *User) CreateUser(context context.Context) (err error) {
	_, err = db.Collection("users").InsertOne(context, user)
	if err != nil {
		glog.Error(err)
		return
	}
	return
}

func LoginUser(context context.Context, l Login) (jwt string, err error) {
	u := User{}
	err = db.Collection("users").FindOne(context, bson.M{"username": l.Username}).Decode(&u)
	if err != nil {
		glog.Error(err)
		return
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password))
	if err != nil {
		glog.Error(err)
		return
	}
	// generate jwt
	jwt = "12345"

	_, err = db.Collection("users").UpdateOne(context, bson.M{"username": u.Username}, bson.M{"$set": bson.M{"jwt": jwt, "loggedon": true, "lastonline": time.Now().Unix()}})
	if err != nil {
		glog.Error(err)
		return
	}
	return
}
