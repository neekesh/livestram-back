package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"bootstrap/webrtc/chat"
	constants "bootstrap/webrtc/constants"
	"bootstrap/webrtc/webrtc"

	"github.com/gin-gonic/gin"
	guid "github.com/google/uuid"
	"github.com/gorilla/websocket"
	pion "github.com/pion/webrtc/v3"
)

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type RoomControllers struct {
	// Room repository.RoomRepository
}

func NewRoomControllers() RoomControllers {
	return RoomControllers{}
}

func (rc RoomControllers) GetAllRoom(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "get all teh Room",
		"data": "Rooms",
	})
}

func (rc RoomControllers) CreateRoom(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/room/%s", guid.New().String()))
}

func (rc RoomControllers) JoinRoom(ctx *gin.Context) {
	uuid := ctx.Param("uid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "unique id is not in the params",
		})
	}
	ws := "ws"
	uuid, suuid, _ := rc.CreateOrGetRoom(uuid)

	ctx.JSON(http.StatusOK, gin.H{
		"room_web_socket_addr":  fmt.Sprintf("%s://%s/rooms/websockets/%s", ws, ctx.Request.Host, uuid),
		"room_link":             fmt.Sprintf("%s://%s/rooms/%s", ctx.Request.URL.Scheme, ctx.Request.Host, uuid),
		"chat_web_socket_addr":  fmt.Sprintf("%s://%s/rooms/websockets/chat/%s", ws, ctx.Request.Host, uuid),
		"viewer_websocket_addr": fmt.Sprintf("%s://%s/rooms/websockets/viewer/%s", ws, ctx.Request.Host, uuid),
		"stream_link":           fmt.Sprintf("%s://%s/stream/%s", ctx.Request.URL.Scheme, ctx.Request.Host, suuid),
		"type":                  "room",
	})
}

func (rc RoomControllers) ChatRoom(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Update Room data",
		"data": "",
	})
}

func (rc RoomControllers) ChatRoomSocket(ctx *gin.Context) {

	uuid := ctx.Param("uuid")
	if uuid != "" {
		return
	}
	webrtc.RooomLock.Lock()

	room := webrtc.Rooms(uuid)

	webrtc.RooomLock.Unlock()
	if room == nil {
		return
	}
	if room.Hub == nil {
		return
	}
	socket := &websocket.Conn{}
	chat.PeerChatConn(socket, room.Conn)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Update Room data",
		"data": "",
	})
}

func (rc RoomControllers) RoomViewerSocket(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		return
	}
	webrtc.RoomsLock.Lock()
	if peer, ok := webrtc.Rooms(uuid); ok {
		webrtc.RoomsLock.Unlock()
		rc.RoomViewerConn(ctx, peer.Peer)
		return
	}
	webrtc.RoomsLock.Unlock()
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"msg":  "Update Room data",
	// 	"data": "",
	// })
}

func (rc RoomControllers) CreateOrGetRoom(uuid string) (string, string, Room) {
	webrtc.RoomsLock.Lock()
	defer webrtc.RoomsLock.Unlock()
	hash := sha256.new()
	hash.Write([]byte(uuid))
	suuid := fmt.Sprintf("%x", hash.Sum(nil))
	if room := webrtc.Rooms[uuid]; room != nil {
		if _, ok := webrtc.Streams[suuid]; !ok {
			webrtc.Stream[suuid] = room
		}
		return uuid, suuid, room
	}
	hub := chat.NewHub()
	peer := &webrtc.Peers{}
	peer.TrackLocals = make(map[string]*pion.TrackLocalStaticRTP)
	room := &webrtc.Rooms{
		Peers: peer,
		Hub:   hub,
	}
	webrtc.Rooms[uuid] = room
	webrtc.Streams[uuid] = room
	go hub.Run()
	return uuis, suuid, room
}

func (rc RoomControllers) RoomSocket(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "unique id is not in the params",
		})
	}
	_, _, room := rc.CreateOrGetRoom(uuid)
	conn, err := constants.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("Received message:", string(message))

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (rc RoomControllers) RoomViewerConn(conn *websocket.Conn, peer *webrtc.Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer conn.Close()
	for {
		select {
		case <-ticker.C:
			websocket, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			websocket.Write([]byte(fmt.Sprintf("%d", len(peer.Connections))))

		}
	}
}

type WebsocketMsg struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
