package labstackecho

import (
	"accounting/application/apperror"
	"accounting/infrastructure/log"
	"accounting/infrastructure/util"
	"accounting/usecase/createjournal"
	"net/http"

	"github.com/labstack/echo/v4"
)

// createJournalHandler ...
func (r *Controller) createJournalHandler(inputPort createjournal.Inport) echo.HandlerFunc {

	return func(c echo.Context) error {

		ctx := log.Context(c.Request().Context(), "createjournal")

		var req createjournal.InportRequest

		if err := c.Bind(&req); err != nil {
			newErr := apperror.FailUnmarshalResponseBodyError
			log.Error(ctx, err.Error())
			return c.JSON(http.StatusBadRequest, util.MustJSON(NewErrorResponse(newErr)))
		}

		log.Info(ctx, util.MustJSON(req))

		res, err := inputPort.Execute(ctx, req)
		if err != nil {
			log.Error(ctx, err.Error())
			return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
		}

		log.Info(ctx, util.MustJSON(res))
		return c.JSON(http.StatusOK, NewSuccessResponse(res))
	}
}
