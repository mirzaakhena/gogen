package nethttp

import (
	"accounting/usecase/createjournal"
	"net/http"
)

type Controller struct {
	Router              *http.ServeMux
	CreateJournalInport createjournal.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.HandleFunc("/createjournal", r.authorized(r.createJournalHandler(http.MethodPost, r.CreateJournalInport)))
}
