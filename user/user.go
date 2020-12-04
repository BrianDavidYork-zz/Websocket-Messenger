package user

import (
	"WebsocketMessenger/db"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Message string
	Data    interface{}
}

func Create(res http.ResponseWriter, req *http.Request) {
	type CreateUser struct {
		Username string
		Password string
	}

	r := Response{}
	u := CreateUser{}

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(r)
		return
	}

	var newUser db.User

	newUser.Username = u.Username

	// hash password
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11)
	newUser.Password = string(hashBytes)

	now := time.Now().Unix()
	newUser.Created = now
	newUser.LastOnline = now

	newUser.LoggedOn = true

	token, err := newUser.CreateUser(req.Context())

	r.Message = "User Created"
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Authorization", "Bearer"+token)
	json.NewEncoder(res).Encode(r)
}

func Login(res http.ResponseWriter, req *http.Request) {
	l := db.Login{}
	r := Response{}

	err := json.NewDecoder(req.Body).Decode(&l)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(r)
		return
	}

	token, err := db.LoginUser(req.Context(), l)
	if err != nil {
		r.Message = "User not found"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(r)
		return
	}

	r.Message = "Login Successful"
	r.Data = token
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	r := Response{}

	username, err := JwtAuthorize(req)
	if err != nil {
		glog.Info(err)
		r.Message = "Error"
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(r)
		return
	}

	err = db.Logout(req.Context(), username)

	r.Message = "Logged Out"
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}

func Profile(res http.ResponseWriter, req *http.Request) {
	r := Response{}

	username, err := JwtAuthorize(req)
	if err != nil {
		glog.Info(err)
		r.Message = "Error"
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(r)
		return
	}
	glog.Info(username)
	u := mux.Vars(req)
	glog.Info(u["username"])

	profile, err := db.GetProfile(req.Context(), u["username"])
	if err != nil {
		glog.Info(err)
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(r)
		return
	}

	r.Message = "Profile Retrieved"
	r.Data = profile
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}

func JwtAuthorize(req *http.Request) (username string, err error) {
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
