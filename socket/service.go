// Serves the WS Service

package socket

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

type WSService struct {
	Hub Hub
}

// Connection handler, upgrades the connection and registers the
// connection with the hub
func (s *WSService) Handler(w http.ResponseWriter, r *http.Request) {
	// Only support GET requests
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	// Upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	// Create Connection instance
	c := &connection{send: make(chan []byte, 256), ws: ws}
	// Register the connection
	s.Hub.register <- c
	// Start the writer for the conneciton
	go c.writer()
}

// Create a new WS Service
func NewWSService(h Hub) WSService {
	return WSService{
		Hub: h,
	}
}
