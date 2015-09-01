// Handle REST calls to POST /events/pause
// This is called when the player enters a pause state and state needs to be
// altered in Redis.

package rest

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	v "github.com/asaskevich/govalidator"
	"github.com/thisissoon/FM-Perceptor/events"
	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
)

// POST /events/pause JSON Request Body
type pauseCreateReqBody struct {
	Start string `json:"start" valid:"iso8601,required"`
}

// POST /events/pause HTTP Handler
func PauseCreateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	rbody := &pauseCreateReqBody{}

	// Decode JSON
	err := decoder.Decode(&rbody)
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Validate
	res, err := v.ValidateStruct(rbody)
	if err != nil {
		log.Debug(res)
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// Set the pause redis keys
	if err := events.SetPauseState(c.Env["REDIS"].(*redis.Client), rbody.Start); err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// We got here! It's alllll good.
	w.WriteHeader(201)
}
