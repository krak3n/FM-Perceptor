// Redis Middleware
// Places a Redis Client instance on the Request Application Context
// so redis can be used in HTTP Handlers

package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
)

// Use c.Env["REDIS"] to access the Redis Client Instance
func RedisClient(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		c.Env["REDIS"] = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		log.Debug("Redis Client Created")

		// Call the next handler
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}
