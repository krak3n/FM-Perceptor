// REST API Handlers for /events

package main

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func volumeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("vol"))
}

// Handle sendiong a mute change event (POST /events/mute)
func muteHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mute"))
}
