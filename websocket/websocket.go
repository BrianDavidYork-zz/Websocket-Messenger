package websocket

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"net"
	"net/http"
)

var connectionTable map[string]*net.Conn

func Create(res http.ResponseWriter, req *http.Request) {
	// using github.com/gobwas/ws package and docs
	glog.Info("Create Websocket")
	conn, _, _, err := ws.UpgradeHTTP(req, res)
	if err != nil {
		glog.Error(err)
	}
	// echo chat from docs
	go func() {
		defer conn.Close()

		for {
			_, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				glog.Error(err)
			}
			err = wsutil.WriteServerMessage(conn, op, []byte("same response all the time!"))
			if err != nil {
				glog.Error(err)
			}
		}
	}()
}

func SendWebsocketMessage() (err error) {
	return
}
