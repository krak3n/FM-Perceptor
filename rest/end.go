// Handle REST calls to POST /events/end

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

// POST /events/end expected request body
type endCreateReqBody struct {
	Track string `json:"uri" valid:"required"`
	User  string `json:"user" valid:"required"`
}

// POST /events/end HTTP Handler
func EndCreateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	rbody := &endCreateReqBody{}

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

	// Publish event
	if err := events.PublishEndEvent(
		c.Env["REDIS"].(*redis.Client),
		rbody.Track,
		rbody.User); err != nil {

		log.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// We got to the end - everything went fine!
	w.WriteHeader(201)
}
