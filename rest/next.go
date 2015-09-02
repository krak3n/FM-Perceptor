// Get the next Track from the Playlist Queue

package rest

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
)

const playlistKey = "fm:player:queue"

// Get the next track from the Redis Playlist queue - if no tracks
// exist in the queue then we will 404
func GetNextTrackHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the redis client
	red := c.Env["REDIS"].(*redis.Client)

	result, err := red.LPop(playlistKey).Result()
	if err == redis.Nil {
		log.Debug("Playlist Empty")
		http.Error(w, http.StatusText(404), 404)
		return
	} else if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// Write the result
	w.Write([]byte(result))
}
