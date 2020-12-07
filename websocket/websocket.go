package websocket

import (
	"WebsocketMessenger/db"
	"bytes"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"net/http"
	"strings"
	"sync"
)

var connectionTable = make(map[string]*net.Conn)

func Create(res http.ResponseWriter, req *http.Request) {
	m := sync.Mutex{}

	// using github.com/gobwas/ws package and docs
	conn, _, _, err := ws.UpgradeHTTP(req, res)
	if err != nil {
		glog.Error(err)
	}

	// on connect, client sends jwt to authorize
	token, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		glog.Error(err)
	}

	username, err := checkJwt(string(token))
	glog.Info(username)
	// if jwt is not authorized, close websocket
	if err != nil || username == "" {
		err = conn.Close()
		if err != nil {
			glog.Info(err)
		}
	} else {
		// if jwt is authorized, save user and connection to table, then listen for messages
		m.Lock()
		connectionTable[username] = &conn
		m.Unlock()
		for {
			msg, _, err := wsutil.ReadClientData(conn)

			// if websocket is closed, remove from map
			if err != nil && strings.Contains(err.Error(), "closed") {
				glog.Info(err)
				_, ok := connectionTable[username]
				if ok {
					m.Lock()
					delete(connectionTable, username)
					m.Unlock()
					break
				}
			} else {
				// on normal read

				type Typing struct {
					User    string
					Message string
					ConvId  string
				}

				t := Typing{}
				reader := bytes.NewReader(msg)

				err := json.NewDecoder(reader).Decode(&t)
				if err != nil {
					glog.Info("Error decoding message")
				}

				Id, err := primitive.ObjectIDFromHex(t.ConvId)
				if err != nil {
					glog.Info("Invalid conversation Id")
				}

				// get conversation
				conv, err := db.GetConversationById(context.Background(), Id)
				if err != nil {
					glog.Error(err)
					return
				}

				// get non-sender chat members
				members := conv.Members

				for _, v := range members {
					if v != t.User {
						// get right conn
						userConn, ok := connectionTable[v]
						if ok {
							u := *userConn
							err = wsutil.WriteServerMessage(u, 1, msg)
							if err != nil {
								glog.Info("Error writing message")
							}
						}
					}
				}
			}
		}
	}
}

func SendWebsocketMessage(message db.Message, messageType string) {
	// get conversation
	conv, err := db.GetConversationById(context.Background(), message.ConversationId)
	if err != nil {
		glog.Error(err)
		return
	}

	// get non-sender chat members
	members := conv.Members

	// for each, get connection from map
	for _, v := range members {
		if v != message.Sender {
			conn, ok := connectionTable[v]
			if ok {
				c := *conn
				err = wsutil.WriteServerMessage(c, 1, []byte(messageType+", "+message.ConversationId.Hex()))
				if err != nil {
					glog.Info("Error sending websocket notification")
				}
			}
		}
	}
	return
}

func checkJwt(bearerToken string) (username string, err error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			glog.Error(err)
		}
		return []byte("WaterCooler123"), nil
	})
	if err != nil {
		glog.Error(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		glog.Error(err)
		return
	}
	username = claims["username"].(string)
	return
}
