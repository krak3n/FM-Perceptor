// Redis Middleware
// Places a Redis Client instance on the Request Application Context
// so redis can be used in HTTP Handlers

package middleware

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
)

// Use c.Env["REDIS"] to access the Redis Client Instance
func RedisClient(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		c.Env["REDIS"] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis_host"), viper.GetString("redis_port")),
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		log.Debugf("Redis @ %s:%s", viper.GetString("redis_host"), viper.GetString("redis_port"))

		// Call the next handler
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}
