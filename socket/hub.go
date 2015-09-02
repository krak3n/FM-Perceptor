// Sets up a Message Hub

package socket

import log "github.com/Sirupsen/logrus"

// Handles sending messages to all connected clients
type Hub struct {
	// Stores all connected clients
	connections map[*connection]bool

	// Messages to broadcast to the connected clients
	Broadcast chan []byte

	// Register a new client
	register chan *connection

	// Unregister clients, removing them from the pool
	unregister chan *connection
}

// Listens on the channels for messages and performs the
// relivant functionality
func (h *Hub) Run() {
	log.Info("Starting Hub")
	for {
		select {
		// On register messages, store the connection
		case c := <-h.register:
			log.Debug("Client Registered: ", c.ws.RemoteAddr())
			h.connections[c] = true
		// On unregister messsages, delete the connection from the pool
		case c := <-h.unregister:
			log.Debug("Client Unregistered: ", c.ws.RemoteAddr())
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		// On incoming messages, loop over connected clients and send
		// the message
		case m := <-h.Broadcast:
			for c := range h.connections {
				log.Debug("Broadcast Message: ", string(m[:]))
				select {
				// Put the message on the connections send channel
				case c.send <- m:
				default:
					log.Debug("Client no longer active: ", c.ws.RemoteAddr())
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

// Create a new message hub
func NewHub() Hub {
	return Hub{
		Broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
	}
}
