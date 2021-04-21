package gochi

import (
	"accounting/usecase/createjournal"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Controller struct {
	Router              *chi.Mux
	CreateJournalInport createjournal.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.Method(http.MethodPost, "/createjournal", r.authorized(r.createJournalHandler(http.MethodPost, r.CreateJournalInport)))
}
