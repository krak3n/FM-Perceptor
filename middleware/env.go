// Sets up a simple Env map to store on the request Context

package middleware

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// Use c.Env["FOO"] = "Bar"
func SetupEnv(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if c.Env == nil {
			c.Env = make(map[interface{}]interface{})
		}

		// Call the next handler
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}
