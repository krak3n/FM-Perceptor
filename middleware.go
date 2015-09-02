// Middleware Methods for REST API Endpoints

package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
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

// HMAC Verification
// Only secific clients are allowed to use the API
func HMACVerification(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		request_data, _ := ioutil.ReadAll(r.Body)
		// Put the data back on the request object
		r.Body = ioutil.NopCloser(bytes.NewBuffer(request_data))
		// Get the request Signature Header
		creds := bytes.SplitN([]byte(r.Header.Get("Signature")), []byte(":"), 2)
		// Ensure we have creds
		if len(creds) != 2 {
			log.Warn("Malformed Request Signature Sent")
			http.Error(w, http.StatusText(400), 400)
			return
		}
		// Assign vars for easy ref
		// client_id := creds[0]
		client_sig := creds[1]
		// Get the Secret Key for the Client
		secret := "foo" // TODO: get from somewhere
		// Encode he request body
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(request_data)
		expected_sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		// Ensure Signatures match
		if string(client_sig[:]) != expected_sig {
			log.Warn("Invalid HMAC Signature - Possible Man in the Middle Attack")
			http.Error(w, http.StatusText(400), 400)
			return
		}

		// Call the next handler on success
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}
