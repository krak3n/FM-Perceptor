// Middleware Methods for REST API Endpoints

package main

import (
	"net/http"

	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
)

// Sets up a Redis Client to be used in Requests - Use c.Env["REDIS"].get("
func RedisClient(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if c.Env == nil {
			c.Env = make(map[interface{}]interface{})
		}

		c.Env["REDIS"] = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}
