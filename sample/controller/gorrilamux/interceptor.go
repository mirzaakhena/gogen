package gorrilamux

import (
	"net/http"
)

// authorized is an interceptor
func (r *Controller) authorized(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		authorized := true

		if !authorized {
			http.Error(w, "Not Authorized", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}