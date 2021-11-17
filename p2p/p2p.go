package p2p

import (
	"fmt"
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
	utils.HandleError(err)
	initPeer(conn, "xx","xx")
}


func AddToPeer(address, port string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws",address,port),nil)
	utils.HandleError(err)
	initPeer(conn, address, port)
}