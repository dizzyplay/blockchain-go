package p2p

import (
	"github.com/dizzyplay/blockchain-go/utils"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}
var conns []*websocket.Conn

func Upgrade(rw http.ResponseWriter, r *http.Request){
	upgrader.CheckOrigin = func (r *http.Request)bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	conns = append(conns, conn)
	utils.HandleError(err)
	for {
		_,p,err := conn.ReadMessage()
		if err != nil {
			break
		}
		for _, aConn := range conns {
			if aConn != conn {
				utils.HandleError(aConn.WriteMessage(websocket.TextMessage, p))
			}
		}
	}
}
