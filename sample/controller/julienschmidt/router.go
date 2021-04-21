package julienschmidt

import (
	"accounting/usecase/createjournal"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Controller struct {
	Router              *httprouter.Router
	CreateJournalInport createjournal.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.Handle(http.MethodPost, "/createjournal", r.authorized(r.createJournalHandler(r.CreateJournalInport)))
}
