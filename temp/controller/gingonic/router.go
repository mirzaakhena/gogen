package gingonic

import (
	"github.com/gin-gonic/gin"
	"github.com/mirzaakhena/gogen2/temp/usecase/createjournal"
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
