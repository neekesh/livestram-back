package webrtc

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	pion "github.com/pion/webrtc/v3"
)

func RoomConnection(conn *websocket.Conn, peers *Peers) {
	var config = pion.Configuration{}
	peerConnection, err := pion.NewPeerConnection(config)
	if err != nil {
		log.Println(err.Error())
		return
	}
	newPeer := PeerConnectionState{
		PeerConnection: peerConnection,
		WebSocket:      &ThreadSafeWriter{},
		Conn:           conn,
		Mutex:          sync.Mutex,
	}
}
