package internal

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 8192
)

// Client represents a single connected client.
//
// The client should die with the connection.
type Client struct {
	hub        *Hub
	connection *websocket.Conn
	send       chan []byte
}

// ReadPump begins accepting messages from the client.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.connection.Close()
	}()

	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.send <- message
	}
}

// WritePump begins sending messages from the hub to the client.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The Hub closed the channel.
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
			}

			w, err := c.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs establishes a websocket connection and begins
// handling messages for it.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:        hub,
		connection: conn,
		send:       make(chan []byte, 8192),
	}
	client.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
