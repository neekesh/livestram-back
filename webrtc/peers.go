package webrtc

import (
	"bootstrap/webrtc/chat"
	"log"
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
	peer.ListLock.Lock()
	defer func() {
		peer.ListLock.Lock()
		peer.SignalPeerConnection()
	}()

	delete(peer.TrackLocals, t.ID())

}

func (peer *Peers) AddTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	peer.ListLock.Lock()

	defer func() {
		peer.ListLock.Unlock()
		peer.SignalPeerConnection()
	}()
	track, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		log.Println("error in track add ", err.Error())
	}
	peer.TrackLocals[t.ID()] = track
	return track
}

func (peer *Peers) SignalPeerConnection() {
	peer.ListLock.Lock()
	defer func() {
		peer.ListLock.Unlock()
		peer.DispatchkeyFrame()
	}()
	attemptSync := func(tryAgain bool) {
		for i := range peer.Connections {
			if peer.Connections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				peer.Connections = append(peer.Connections[:i], peer.Connections[i+1]...)
				log.Println("a", peer.Connections)
				return true
			}
			existingSender := map[string]bool{}
			for _, sender := range peer.Connections[i].PeerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}
				existingSender[sender.Track().ID()] = true
				if _, ok := peer.TrackLocals[sender.Track().ID()]; !ok {
					if err := peer.Connections[i].PeerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}
			for _, reciever := range peer.Connections[i].PeerConnection.GetReceivers() {
				if reciever.Track() == nil {
					continue
				}
				existingSender[reciever.Track().ID()] = true
			}
			for trackID := range peer.TrackLocals {
				if _, ok := existingSender[trackID]; !ok {
					if _, err := peer.Connections[i].PeerConnection.AddTrack(peer.TrackLocals[trackID]); err != nil {
						return nil
					}
				}
			}
		}
	}
}
