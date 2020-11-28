package main

import (
	"context"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"time"
)

func main() {
	// consider using env variables

	startMongo()

	router := mux.NewRouter()

	// API ROUTES

	// user
	//router.HandleFunc("user", user.Create).Methods("POST")
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

	// jwt auth middleware

	// start server
	glog.Info("Starting messenger api on port 7000")
	if err := http.ListenAndServe(":7000", router); err != nil {
		glog.Fatal("ListenAndServe: ", err)
	}
}

func startMongo() {
	// The code in this function comes from the go mongo-driver docs
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	glog.Info("MongoDb started")
}
