package conversation

import (
	"WebsocketMessenger/db"
	"WebsocketMessenger/response"
	"encoding/json"
	"github.com/golang/glog"
	"net/http"
	"time"
)

func Create(res http.ResponseWriter, req *http.Request) {
	r := response.Response{}

	u := req.Context().Value("username")
	username := u.(string)

	type Recipient struct {
		Recipient string
	}

	rcp := Recipient{}

	err := json.NewDecoder(req.Body).Decode(&rcp)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	c := db.Conversation{}
	c.Members = append(c.Members, username, rcp.Recipient)
	c.Created = time.Now().Unix()
	chatId, err := c.Create(req.Context(), username, rcp.Recipient)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	r.Message = "Conversation Created"
	r.Data = chatId
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(r)
	if err != nil {
		glog.Info(err)
	}
}

func GetAllConversations(res http.ResponseWriter, req *http.Request) {
	r := response.Response{}

	u := req.Context().Value("username")
	username := u.(string)

	conversations, err := db.GetConversations(req.Context(), username)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	r.Message = "Conversations Retrieved"
	r.Data = conversations
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(r)
	if err != nil {
		glog.Info(err)
	}
}
