package websocket

import (
	"github.com/gobwas/ws"
	"github.com/golang/glog"
	"net"
	"net/http"
)

func Create(res http.ResponseWriter, req *http.Request) {
	ln, err := net.Listen("tcp", ":7000")
	if err != nil {
		glog.Info(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		glog.Info(err)
	}

	handshake, err := ws.Upgrade(conn)
	if err != nil {
		glog.Info(err)
	}
	glog.Info(handshake)
}
