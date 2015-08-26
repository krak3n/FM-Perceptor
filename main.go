// Serves the Websocket Server

package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

// Upgrade instance to upgrade the connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Connection handler, upgrades the connection and registers the
// connection with the hub
func serve(w http.ResponseWriter, r *http.Request) {
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
	go h.run()

	// Redis Connection
	s := NewSubscription()
	go s.consume()

	log.Debug("Starting Websocket Server on :9000")
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
