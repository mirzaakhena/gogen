package labstackecho

import (
	"accounting/usecase/createjournal"
	"github.com/labstack/echo/v4"
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
