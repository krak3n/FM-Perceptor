// Handle Volume Change Events
// Called when the players alters its Volume from a user drivem event

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

// PUT /volume Request Body
type volumeUpdateReqBody struct {
	Level int `json:"level" valid:"required"`
}

// PUT /volume HTTP Handler
func VolumeUpdateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	rbody := &volumeUpdateReqBody{}

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
	if err := events.PublishVolumeEvent(c.Env["REDIS"].(*redis.Client), rbody.Level); err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// We got here! It's alllll good.
	w.WriteHeader(200)
}
