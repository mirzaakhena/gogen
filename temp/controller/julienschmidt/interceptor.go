package julienschmidt

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (r *Controller) authorized(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		authorized := true

		if !authorized {
			http.Error(w, "Not Authorized", http.StatusForbidden)
			return
		}

	}
}
