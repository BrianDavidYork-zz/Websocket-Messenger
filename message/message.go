package message

import (
	"WebsocketMessenger/db"
	"WebsocketMessenger/user"
	"encoding/json"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Response struct {
	Message string
	Data    interface{}
}

func Create(res http.ResponseWriter, req *http.Request) {
	r := Response{}
	m := db.Message{}

	username, err := user.JwtAuthorize(req)
	if err != nil {
		glog.Info(err)
		r.Message = "Error"
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(r)
		return
	}

	err = json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(r)
		return
	}

	m.Sender = username
	m.State = 0
	mid, err := m.Create(req.Context())
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(r)
		return
	}

	// send ws notification to other member of conv

	r.Message = "Message Created"
	r.Data = mid
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}

func Edit(res http.ResponseWriter, req *http.Request) {
	type EditMessage struct {
		MessageId string
		Message   string
	}

	r := Response{}
	e := EditMessage{}

	_, err := user.JwtAuthorize(req)
	if err != nil {
		glog.Info(err)
		r.Message = "Error"
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(r)
		return
	}

	err = json.NewDecoder(req.Body).Decode(&e)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(r)
		return
	}

	chatId, err := primitive.ObjectIDFromHex(e.MessageId)
	if err != nil {
		r.Message = "Invalid Chat Id"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(r)
		return
	}

	err = db.EditMessage(req.Context(), chatId, e.Message)
	if err != nil {
		r.Message = "Error Editing Message"
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(r)
		return
	}

	// send ws notification to other member of conv

	r.Message = "Message Edited"
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}
