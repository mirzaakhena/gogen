package julienschmidt

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (r *Controller) authorized(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		authorized := true

		if !authorized {
			http.Error(w, "Not Authorized", http.StatusForbidden)
			return
		}

		next(w, r, ps)

	}
}
