// Adds Extra Validtors to Go Validator

package middleware

import (
	"net/http"
	"time"

	v "github.com/asaskevich/govalidator"
	"github.com/zenazn/goji/web"
)

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
