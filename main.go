package main

import (
	"./db"
	"./user"
	"./websocket"
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
	router.HandleFunc("/user/{username}", user.Profile).Methods("GET")
	router.HandleFunc("/user/login", user.Login).Methods("POST")
	router.HandleFunc("/user/logout", user.Logout).Methods("POST")

	// websocket
	router.HandleFunc("/websocket", websocket.Create).Methods("GET")

	// conversation
	//router.HandleFunc("/conversation", conversation.Create).Methods("POST")
	//router.HandleFunc("/conversation", conversation.GetAllConversations).Methods("GET")

	// messages
	//router.HandleFunc("/message", message.Create).Methods("POST")
	//router.HandleFunc("/message", message.Edit).Methods("PUT")
	//router.HandleFunc("/message/{id}", message.Delete).Methods("DELETE")
	//router.HandleFunc("/message", message.GetMessages).Methods("GET")

	// start server
	glog.Info("Starting messenger api on port 7000")
	if err := http.ListenAndServe(":7000", router); err != nil {
		glog.Fatal("ListenAndServe: ", err)
	}
}

//func jwtMiddleware(req *http.Request) (err error) {
//	header := req.Header.Get("Authorization")
//}

// TODO

// jwt authorization - turn into middleware
// env variables - jwt secret, mongo url, api port
// log errors
