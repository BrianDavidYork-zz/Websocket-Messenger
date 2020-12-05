package conversation

import (
	"WebsocketMessenger/db"
	"WebsocketMessenger/response"
	"encoding/json"
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
		json.NewEncoder(res).Encode(r)
		return
	}

	c := db.Conversation{}
	c.Members = append(c.Members, username, rcp.Recipient)
	c.Created = time.Now().Unix()
	chatId, err := c.Create(req.Context(), username, rcp.Recipient)
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
	r := response.Response{}

	u := req.Context().Value("username")
	username := u.(string)

	conversations, err := db.GetConversations(req.Context(), username)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(r)
		return
	}

	// send ws notification to recipient

	r.Message = "Conversations Retrieved"
	r.Data = conversations
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(r)
}
