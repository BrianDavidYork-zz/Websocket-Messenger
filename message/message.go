package message

import (
	"WebsocketMessenger/db"
	"WebsocketMessenger/response"
	"WebsocketMessenger/websocket"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func Create(res http.ResponseWriter, req *http.Request) {
	r := response.Response{}
	m := db.Message{}

	u := req.Context().Value("username")
	username := u.(string)

	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	m.Sender = username
	m.State = 0
	m.Created = time.Now().Unix()
	mid, err := m.Create(req.Context())
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	conv, err := db.GetConversationById(req.Context(), m.ConversationId)
	if err != nil {
		glog.Info(err)
	}

	// mark new messages for all non-requesting users in conversation
	for _, v := range conv.Members {
		if v != username {
			// err = db.MarkConversationUnseen(req.Context(), v, m.ConversationId)
			if err != nil {
				glog.Info(err)
			}
		}
	}

	// send ws notification to other member of conv
	websocket.SendWebsocketMessage(m, "New Message")

	r.Message = "Message Created"
	r.Data = mid
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(r)
	if err != nil {
		glog.Info(err)
	}
}

func Edit(res http.ResponseWriter, req *http.Request) {
	type EditMessage struct {
		MessageId string
		Message   string
	}

	r := response.Response{}
	e := EditMessage{}

	u := req.Context().Value("username")
	username := u.(string)

	err := json.NewDecoder(req.Body).Decode(&e)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	messageId, err := primitive.ObjectIDFromHex(e.MessageId)
	if err != nil {
		r.Message = "Invalid Message Id"
		res.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// get message by Id
	msg, err := db.GetMessageById(req.Context(), messageId)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// compare username and message.Sender
	if msg.Sender != username {
		r.Message = "Not Authorized to Edit Message"
		res.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	err = db.EditMessage(req.Context(), messageId, e.Message)
	if err != nil {
		r.Message = "Error Editing Message"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// send ws notification to other member of conv
	websocket.SendWebsocketMessage(msg, "Message Edited")

	r.Message = "Message Edited"
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(r)
	if err != nil {
		glog.Info(err)
	}
}

func Delete(res http.ResponseWriter, req *http.Request) {
	r := response.Response{}

	u := req.Context().Value("username")
	username := u.(string)

	vars := mux.Vars(req)
	msgId := vars["id"]
	if msgId == "" {
		r.Message = "ID Parameter Required"
		res.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	messageId, err := primitive.ObjectIDFromHex(msgId)
	if err != nil {
		r.Message = "Invalid Message Id"
		res.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// get message by Id
	msg, err := db.GetMessageById(req.Context(), messageId)
	if err != nil {
		r.Message = "Error"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// compare username and message.Sender
	if msg.Sender != username {
		r.Message = "Not Authorized to Delete Message"
		res.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	err = db.DeleteMessage(req.Context(), messageId)
	if err != nil {
		r.Message = "Error Deleting Message"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// send ws notification to other member of conv
	websocket.SendWebsocketMessage(msg, "Message Deleted")

	r.Message = "Message Deleted"
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(r)
	if err != nil {
		glog.Info(err)
	}
}

func GetMessages(res http.ResponseWriter, req *http.Request) {
	r := response.Response{}

	u := req.Context().Value("username")
	username := u.(string)

	// get id from route param
	vars := mux.Vars(req)
	convId := vars["id"]
	if convId == "" {
		r.Message = "ID Parameter Required"
		res.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	conversationId, err := primitive.ObjectIDFromHex(convId)
	if err != nil {
		r.Message = "Invalid Message Id"
		res.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// get conversation by Id
	conversation, err := db.GetConversationById(req.Context(), conversationId)
	if err != nil {
		r.Message = "Error Retrieving Conversation"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// make sure requesting user is a conversation member
	var permitted bool
	permitted = false

	for _, v := range conversation.Members {
		if v == username {
			permitted = true
		}
	}

	if !permitted {
		r.Message = "Not Authorized to Retrieve Messages"
		res.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	messages, err := db.GetMessagesByConversation(req.Context(), conversationId)
	if err != nil {
		r.Message = "Error Retrieving Messages"
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			glog.Info(err)
		}
		return
	}

	// mark conversation as seen for requesting user
	// err = db.MarkConversationSeen(req.Context(), username, conversationId)
	// if err != nil {
	// glog.Info(err)
	// }

	r.Message = "Messages Retrieved"
	r.Data = messages
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(r)
	if err != nil {
		glog.Info(err)
	}
}
