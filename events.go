// REST API Handlers for /events

package main

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// Handle sending a play event (POST /events/play)
func playHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("play"))
}

// Handle sending a end event (POST /events/end)
func endHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("end"))
}

// Handle sendiong a volume change event (POST /events/volume)
func volumeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("vol"))
}

// Handle sendiong a mute change event (POST /events/mute)
func muteHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mute"))
}
