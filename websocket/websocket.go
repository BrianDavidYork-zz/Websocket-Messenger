package websocket

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"net"
	"net/http"
)

var connectionTable = make(map[string]*net.Conn)

func Create(res http.ResponseWriter, req *http.Request) {
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
	if err != nil || username == "" {
		conn.Close()
	} else {
		// if jwt is authorized
		// save user and connection to table
		connectionTable[username] = &conn
	}
}

func SendWebsocketMessage(username string, message interface{}) (err error) {
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
