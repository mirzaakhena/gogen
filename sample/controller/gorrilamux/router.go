package gorrilamux

import (
	"accounting/usecase/createjournal"
	"github.com/gorilla/mux"
	"net/http"
)

type Controller struct {
	Router              *mux.Route
	CreateJournalInport createjournal.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.Methods(http.MethodPost).Path("/createjournal").HandlerFunc(r.authorized(r.createJournalHandler(r.CreateJournalInport)))
}
