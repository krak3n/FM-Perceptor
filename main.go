// Serves the Websocket Server

package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

// Upgrade instance to upgrade the connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Connection handler, upgrades the connection and registers the
// connection with the hub
func serveWS(w http.ResponseWriter, r *http.Request) {
	// Only support GET requests
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	// Upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create Connection instance
	c := &connection{send: make(chan []byte, 256), ws: ws}
	// Register the connection
	h.register <- c
	// Start the writer for the conneciton
	go c.writer()
}

// Entrypoint - Runs the WS Server
func main() {
	log.SetLevel(log.DebugLevel)
	// Websocket Message Hub
	go h.run()

	// Redis Connection
	s := NewSubscription()
	go s.consume()

	// Serve the WS Server
	log.Debug("Starting Websocket Server on :9000")

	c := web.New()
	c.Use(RedisClient)

	// WS Connections
	c.Get("/", serveWS)

	// Event REST endpoints
	c.Post("/events/play", playHandler)
	c.Post("/events/end", endHandler)
	c.Post("/events/volume", volumeHandler)
	c.Post("/events/mute", muteHandler)

	graceful.ListenAndServe(":9000", c)
}
