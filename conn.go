// Handles WS Connections as well as Sending / Recieving Messages
// from conected clients

package main

import (
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
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
	maxMessageSize = 512
)

// Represents a WS connecton to Perceptor
type connection struct {
	// The actual WS connection
	ws *websocket.Conn

	// A channel to send messages to the connection
	send chan []byte
}

// Returns the remote address of the WS Connection
func (c *connection) addr() net.Addr {
	return c.ws.RemoteAddr()
}

// Writes a given message type and payload to the WS connection
func (c *connection) write(mt int, payload []byte) error {
	log.Debug("Write Message to: ", c.addr())
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// Writes messages to the WS Connection
func (c *connection) writer() {
	// Create a ticker that will ping the client
	ticker := time.NewTicker(pingPeriod)

	// Ensure we stop the ticker and close the connection when we exit
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		// Get a message from the connections send channel
		case m, ok := <-c.send:
			log.Info("Send Message to: ", c.addr())
			// Not OK, send a close message
			if !ok {
				log.Error("Not Ok")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			// Attempt to write the message to the connection, catching errors
			if err := c.write(websocket.TextMessage, m); err != nil {
				log.Error(err)
				return // do nothing
			}
		case <-ticker.C:
			// Ping the client to keep the connection open
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Error(err)
				return // do nothing
			}
		}
	}
}
