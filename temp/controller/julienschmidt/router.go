package julienschmidt

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mirzaakhena/gogen2/temp/usecase/createjournal"
	"net/http"
)

type Controller struct {
	Router              *httprouter.Router
	CreateJournalInport createjournal.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.Handle(http.MethodPost, "/createjournal", r.authorized(r.createJournalHandler(http.MethodPost, r.CreateJournalInport)))
}
