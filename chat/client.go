package chat

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

var (
	newline = []byte("\n")
	space   = []byte(" ")
)

func (client *Client) ReadPump() {
	defer func() {
		client.Hub.register <- client
		client.Conn.Close()
	}()
	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(
		func(string) error {
			client.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error %v", err)
			}
			log.Printf("error %v", err)
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		client.Hub.broadcast <- message

	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}
			writer, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(message)
			for index := 0; index <= len(client.Send); index++ {
				writer.Write(newline)
				writer.Write(<-client.Send)
			}
			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func PeerChatConn(conn *websocket.Conn, hub *Hub) {
	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	client.Hub.register <- client
	go client.WritePump()
	client.ReadPump()
}
