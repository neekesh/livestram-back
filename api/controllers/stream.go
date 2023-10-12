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

type StreamControllers struct {
	// Stream repository.StreamRepository
}

func NewStreamControllers() StreamControllers {
	return StreamControllers{}
}

func (rc StreamControllers) GetAllStream(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "get all teh Stream",
		"data": "Streams",
	})
}

func (rc StreamControllers) CreateStream(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/room/%s", guid.New().String()))
}

func (rc StreamControllers) JoinStream(ctx *gin.Context) {
	suuid := ctx.Param("suuid")
	if suuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "unique id is not in the params",
		})
	}
	ws := "ws"
	webrtc.RoomsLock.Lock()
	if _, ok := webrtc.Stream[suuid]; ok {
		webrtc.RoomsLock.Unlock()
		ctx.JSON(http.StatusOK, gin.H{
			"stream_web_socket_addr": fmt.Sprintf("%s://%s/stream/websockets/%s", ws, ctx.Request.Host, suuid),
			// "room_link":             fmt.Sprintf("%s://%s/rooms/%s", ctx.Request.URL.Scheme, ctx.Request.Host, uuid),
			"chat_web_socket_addr":  fmt.Sprintf("%s://%s/stream/websockets/chat/%s", ws, ctx.Request.Host, suuid),
			"viewer_websocket_addr": fmt.Sprintf("%s://%s/stream/websockets/viewer/%s", ws, ctx.Request.Host, suuid),
			// "stream_link":           fmt.Sprintf("%s://%s/stream/%s", ctx.Request.URL.Scheme, ctx.Request.Host, s  uuid),
			"type": "stream",
		})
		return
	}
	webrtc.RoomsLock.Unlock()
	ctx.JSON(http.StatusBadRequest, gin.H{
		"stream": "false",
		"leave":  "true",
	})

	// uuid, suuid, _ := rc.CreateOrGetStream(uuid)
	//
}

func (rc StreamControllers) ChatStream(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Update Stream data",
		"data": "",
	})
}

func (rc StreamControllers) ChatStreamSocket(ctx *gin.Context) {

	suuid := ctx.Param("suuid")
	if suuid != "" {
		return
	}
	webrtc.RooomLock.Lock()

	// room := webrtc.Streams(suuid)
	if stream, ok := webrtc.Streams[suuid]; ok {
		webrtc.Rooms.Unlock()
		webrtc.StreamConn(ctx, stream.Peers)
		retrun
	}

	webrtc.RooomLock.Unlock()
}

func (rc StreamControllers) StreamViewerSocket(ctx *gin.Context) {
	suuid := ctx.Param("suuid")
	if suuid != "" {
		return
	}
	webrtc.RooomLock.Lock()
	if stream, ok := webrtc.Streams[suuid]; ok {
		webrtc.Rooms.Unlock()
		rc.StreamViewerConn(ctx, stream.Peers)
		retrun
	}
	webrtc.StreamsLock.Unlock()
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"msg":  "Update Stream data",
	// 	"data": "",
	// })
}

func (rc StreamControllers) CreateOrGetStream(uuid string) (string, string, Stream) {
	webrtc.StreamsLock.Lock()
	defer webrtc.StreamsLock.Unlock()
	hash := sha256.new()
	hash.Write([]byte(uuid))
	suuid := fmt.Sprintf("%x", hash.Sum(nil))
	if room := webrtc.Streams[uuid]; room != nil {
		if _, ok := webrtc.Streams[suuid]; !ok {
			webrtc.Stream[suuid] = room
		}
		return uuid, suuid, room
	}
	hub := chat.NewHub()
	peer := &webrtc.Peers{}
	peer.TrackLocals = make(map[string]*pion.TrackLocalStaticRTP)
	room := &webrtc.Streams{
		Peers: peer,
		Hub:   hub,
	}
	webrtc.Streams[uuid] = room
	webrtc.Streams[uuid] = room
	go hub.Run()
	return uuis, suuid, room
}

func (rc StreamControllers) StreamSocket(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "unique id is not in the params",
		})
	}
	_, _, room := rc.CreateOrGetStream(uuid)
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

func (rc StreamControllers) StreamViewerConn(conn *websocket.Conn, peer *webrtc.Peers) {
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
