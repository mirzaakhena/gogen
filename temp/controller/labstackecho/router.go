package labstackecho

import (
	"github.com/labstack/echo/v4"
	"github.com/mirzaakhena/gogen2/temp/usecase/createjournal"
	"net/http"
)

type Controller struct {
	Router              *echo.Router
	CreateJournalInport createjournal.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.Add(http.MethodPost, "/createjournal", r.authorized(r.createJournalHandler(r.CreateJournalInport)))
}
