package main

import (
	"WebsocketMessenger/conversation"
	"WebsocketMessenger/db"
	"WebsocketMessenger/message"
	"WebsocketMessenger/response"
	"WebsocketMessenger/user"
	"WebsocketMessenger/websocket"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func init() {
	flag.Set("logtostderr", "true")
	flag.Parse()
}

func main() {
	// consider using env variables

	db.StartMongo()

	router := mux.NewRouter()

	// API ROUTES
	api := router.PathPrefix("/").Subrouter()

	// middleware
	api.Use(jwtMiddleware)

	// user
	router.HandleFunc("/user", user.Create).Methods("POST")
	api.HandleFunc("/user/{username}", user.Profile).Methods("GET")
	router.HandleFunc("/user/login", user.Login).Methods("POST")
	api.HandleFunc("/user/logout", user.Logout).Methods("POST")

	// conversation
	api.HandleFunc("/conversation", conversation.Create).Methods("POST")
	api.HandleFunc("/conversation", conversation.GetAllConversations).Methods("GET")

	// message
	api.HandleFunc("/message", message.Create).Methods("POST")
	api.HandleFunc("/message", message.Edit).Methods("PUT")
	api.HandleFunc("/message/{id}", message.Delete).Methods("DELETE")
	api.HandleFunc("/message/{id}", message.GetMessages).Methods("GET")

	// websocket (no middleware)
	router.HandleFunc("/websocket", websocket.Create).Methods("GET")

	// start server
	glog.Info("Starting messenger api on port 7000")
	if err := http.ListenAndServe(":7000", router); err != nil {
		glog.Fatal("ListenAndServe: ", err)
	}
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// no jwt needed for these routes
		if req.URL.Path == "/user/login" ||
			req.URL.Path == "/user" {
			next.ServeHTTP(res, req)
		}
		var bearerToken string
		tok := req.Header.Get("Authorization")
		strArr := strings.Split(tok, " ")
		if len(strArr) == 2 {
			bearerToken = strArr[1]
		} else {
			notAuthorized(res)
			return
		}
		token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid jwt")
			}
			return []byte("WaterCooler123"), nil
		})
		if err != nil {
			notAuthorized(res)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			notAuthorized(res)
			return
		}
		username := claims["username"].(string)
		newContext := context.WithValue(req.Context(), "username", username)
		next.ServeHTTP(res, req.WithContext(newContext))
	})
}

func notAuthorized(res http.ResponseWriter) {
	r := response.Response{}
	r.Message = "Invalid Jwt"
	r.Data = nil
	json.NewEncoder(res).Encode(r)
}

// TODO

// FUTURE IMPROVEMENTS
// pagination for messages GET
// multi user conversations

// env variables - jwt secret, mongo url, api port
// check all error messages being returned
// add comments
// finish readme

// websocket - rapid open close bug;
