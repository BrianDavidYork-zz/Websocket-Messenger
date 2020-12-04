package conversation

import (
	"WebsocketMessenger/db"
	"WebsocketMessenger/user"
	"encoding/json"
	"github.com/golang/glog"
	"net/http"
)

type Response struct {
	Message string
	Data    interface{}
}

func Create(res http.ResponseWriter, req *http.Request) {
	r := Response{}

	username, err := user.JwtAuthorize(req)
	if err != nil {
		glog.Info(err)
		r.Message = "Error"
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(r)
		return
	}

	type Username struct {
		Recipient string
	}

	u := Username{}

	err = json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(r)
		return
	}

	c := db.Conversation{}
	c.Members = append(c.Members, username, u.Recipient)
	chatId, err := c.Create(req.Context(), username, u.Recipient)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(r)
		return
	}

	// send ws notification to recipient

	r.Message = "Conversation Created"
	r.Data = chatId
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}

func GetAllConversations(res http.ResponseWriter, req *http.Request) {
	r := Response{}

	username, err := user.JwtAuthorize(req)
	if err != nil {
		glog.Info(err)
		r.Message = "Error"
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(r)
		return
	}

	//conversations, err := db.GetConversations(req.Context(), username)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(r)
		return
	}

	// send ws notification to recipient

	r.Message = "Conversations Retrieved"
	//r.Data = conversations
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}