// Middleware Methods for REST API Endpoints

package main

import (
	"net/http"
	"time"

	v "github.com/asaskevich/govalidator"
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

// Add custom validators to GoValidator
func CustomValidators(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// ISO 8601 Validator
		v.TagMap["iso8601"] = v.Validator(func(str string) bool {
			_, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return false
			}

			return true
		})

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}
