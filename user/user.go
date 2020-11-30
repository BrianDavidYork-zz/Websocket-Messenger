package user

import (
	"../db"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func Create(res http.ResponseWriter, req *http.Request) {
	type CreateUser struct {
		Username string
		Password string
	}

	u := CreateUser{}

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	glog.Info(u.Username)
	glog.Info(u.Password)

	var newUser db.User

	newUser.Username = u.Username

	// hash password
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11)
	newUser.Password = string(hashBytes)

	now := time.Now().Unix()
	newUser.Created = now
	newUser.LastOnline = now

	// generate and save jwt

	newUser.LoggedOn = true

	token, err := newUser.CreateUser(req.Context())

	// send back res
	res.Header().Set("Authorization", "Bearer"+token)
}

func Login(res http.ResponseWriter, req *http.Request) {
	l := db.Login{}

	err := json.NewDecoder(req.Body).Decode(&l)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := db.LoginUser(req.Context(), l)

	res.Header().Set("Authorization", "Bearer "+token)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	username, _ := jwtAuthorize(req)
}

func jwtAuthorize(req *http.Request) (username string, err error) {
	var bearerToken string
	tok := req.Header.Get("Authorization")
	strArr := strings.Split(tok, " ")
	if len(strArr) == 2 {
		bearerToken = strArr[1]
	}
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			glog.Error(err)
		}
		return []byte("WaterCooler123"), nil
	})
	if err != nil {
		glog.Error(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		glog.Error(err)
		return
	}
	username = claims["username"].(string)
	return
}
