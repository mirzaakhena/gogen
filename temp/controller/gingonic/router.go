package gingonic

import (
	"accounting/usecase/createjournal"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	Router              gin.IRouter
	CreateJournalInport createjournal.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.Handle(http.MethodPost, "/createjournal", r.authorized(), r.createJournalHandler(r.CreateJournalInport))
}
