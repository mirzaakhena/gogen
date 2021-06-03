package labstackecho

import (
	"accounting/usecase/createjournal"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Router              *echo.Echo
	CreateJournalInport createjournal.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.Add(http.MethodPost, "/createjournal", r.authorized(r.createJournalHandler(r.CreateJournalInport)))
}
