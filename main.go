package main

import (
	"./db"
	"./user"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	// consider using env variables

	db.StartMongo()

	router := mux.NewRouter()

	// API ROUTES

	// user
	router.HandleFunc("/user", user.Create).Methods("POST")
	//router.HandleFunc("user", user.Edit).Methods("PUT")
	//router.HandleFunc("user/{id}", user.Profile).Methods("GET")
	//router.HandleFunc("user/login", user.Login).Methods("POST")
	//router.HandleFunc("user/logout", user.Login).Methods("POST")

	// websocket
	//router.HandleFunc("websocket", websocket.Create).Methods("POST")

	// conversation
	//router.HandleFunc("conversation", conversation.Create).Methods("POST")
	//router.HandleFunc("conversation/{id}", conversation.GetAllConversations).Methods("GET")

	// messages
	//router.HandleFunc("message", message.Create).Methods("POST")
	//router.HandleFunc("message", message.Edit).Methods("PUT")
	//router.HandleFunc("message", message.Delete).Methods("DELETE")
	//router.HandleFunc("message/{id}", message.GetMessages).Methods("GET")

	// start server
	glog.Info("Starting messenger api on port 7000")
	if err := http.ListenAndServe(":7000", router); err != nil {
		glog.Fatal("ListenAndServe: ", err)
	}
}

//func jwtMiddleware(req *http.Request) (err error) {
//	header := req.Header.Get("Authorization")
//}
