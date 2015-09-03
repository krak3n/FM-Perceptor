// Handle Mute REST Calls

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

// PUT /mute Request Body
type muteUpdateReqBody struct {
	Active bool `json:"active" valid:"required"`
}

// PUT /mute HTTP Handler
func MuteUpdateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	rbody := &muteUpdateReqBody{}

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

	// Set the vol redis keys
	if err := events.PublishMuteEvent(c.Env["REDIS"].(*redis.Client), rbody.Active); err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// We got here! It's alllll good.
	w.WriteHeader(200)
}
