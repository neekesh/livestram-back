package webrtc

import (
	"bootstrap/webrtc/chat"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	pion "github.com/pion/webrtc/v3"
)

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Rooms struct {
	Peers *Peers
	Hub   *chat.Hub
}

type Peers struct {
	ListLock    sync.RWMutex
	Connections []PeerConnectionState
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type PeerConnectionState struct {
	PeerConnection *pion.PeerConnection
	Websocket      *ThreadSafeWriter
}

type ThreadSafeWriter struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

func (tsw *ThreadSafeWriter) WriteJSON(v interface{}) error {
	tsw.Mutex.Lock()
	defer tsw.Mutex.Unlock()
	return tsw.Conn.WriteJSON(v)
}

func (peer *Peers) RemoveTrack(t *webrtc.TrackLocalStaticRTP) {

}

func (peer *Peers) AddTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {

}

func (peer *Peers) SignalPeerConnection() {

}
