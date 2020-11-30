package db

import (
	"context"
	"github.com/dgrijalva/jwt-go"
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

func (user *User) CreateUser(context context.Context) (token string, err error) {
	token, err = generateJwt(*user)
	user.Jwt = token
	_, err = db.Collection("users").InsertOne(context, user)
	if err != nil {
		glog.Error(err)
		return
	}
	return
}

func LoginUser(context context.Context, l Login) (token string, err error) {
	u := User{}
	err = db.Collection("users").FindOne(context, bson.M{"username": l.Username}).Decode(&u)
	if err != nil {
		glog.Error(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password))
	if err != nil {
		glog.Error(err)
		return
	}

	token, err = generateJwt(u)
	if err != nil {
		glog.Error(err)
		return
	}

	_, err = db.Collection("users").UpdateOne(context,
		bson.M{"username": u.Username},
		bson.M{"$set": bson.M{"jwt": token, "loggedon": true, "lastonline": time.Now().Unix()}})
	if err != nil {
		glog.Error(err)
		return
	}
	return
}

func Logout(context context.Context, username string) (err error) {
	glog.Info(username)
	_, err = db.Collection("users").UpdateOne(context,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"jwt": "", "loggedon": false, "lastonline": time.Now().Unix()}})
	if err != nil {
		glog.Error(err)
		return
	}
	return
}

func generateJwt(u User) (token string, err error) {
	claims := jwt.MapClaims{
		"username": u.Username,
		"exp":      time.Now().Add(time.Hour * 12),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString([]byte("WaterCooler123"))
	if err != nil {
		glog.Error(err)
		return
	}
	return
}
