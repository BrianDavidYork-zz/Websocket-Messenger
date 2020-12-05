package user

import (
	"WebsocketMessenger/db"
	"WebsocketMessenger/response"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Create(res http.ResponseWriter, req *http.Request) {
	type CreateUser struct {
		Username string
		Password string
	}

	r := response.Response{}
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
	r := response.Response{}

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
	r := response.Response{}

	u := req.Context().Value("username")
	username := u.(string)

	_ = db.Logout(req.Context(), username)

	r.Message = "Logged Out"
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}

func Profile(res http.ResponseWriter, req *http.Request) {
	r := response.Response{}

	u := mux.Vars(req)

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
