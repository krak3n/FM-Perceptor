// Resume Event Handler

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

// POST /events/resume JSON Request Body
type resumeCreateReqBody struct {
	Durration string `json:"start" valid:"int,required"`
}

// POST /player/resume HTTP Handler
// The player will send the durration the player was paused for - this is then
// used to calculate time remaining of the current song
func ResumeCreateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	rbody := &resumeCreateReqBody{}

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

	// Set the resume redis keys
	if err := events.SetResumeState(c.Env["REDIS"].(*redis.Client), rbody.Durration); err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// We got here! It's alllll good.
	w.WriteHeader(201)
}
