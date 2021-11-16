package p2p

import (
	"github.com/dizzyplay/blockchain-go/utils"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request){
	upgrader.CheckOrigin = func (r *http.Request)bool {
		return true
	}
	_, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)
}