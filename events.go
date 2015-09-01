// REST API Handlers for /events

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	v "github.com/asaskevich/govalidator"
	"github.com/thisissoon/FM-Perceptor/events"
	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
)

type playEvent struct {
	Start string `json:"start" valid:"iso8601,required"`
	Track string `json:"uri" valid:"required"`
	User  string `json:"user" valid:"required"`
}

// Handle sending a play event (POST /events/play)
func playHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	event := &playEvent{}

	// Decode JSON
	err := decoder.Decode(&event)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Validate
	res, err := v.ValidateStruct(event)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(422), 422)
		fmt.Println(res)
		return
	}

	// Publish event
	if err := events.PublishPlayEvent(
		c.Env["REDIS"].(*redis.Client),
		event.Track,
		event.User,
		event.Start); err != nil {

		http.Error(w, http.StatusText(500), 500)
		return
	}

	// We got to the end - everything went fine!
	w.WriteHeader(201)
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
