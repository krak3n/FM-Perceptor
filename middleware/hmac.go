// Provides HMAC Request Verification Via a Signature Header

package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

// Ensures only certain requests are valid to the service
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
