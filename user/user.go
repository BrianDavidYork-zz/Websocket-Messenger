package user

import (
	"../db"
	"encoding/json"
	"github.com/golang/glog"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

	err = newUser.CreateUser(req.Context())

	// send back res
}
