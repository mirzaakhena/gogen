package labstackecho

import (
	"github.com/labstack/echo/v4"
	"github.com/mirzaakhena/gogen2/temp/apperror"
	"github.com/mirzaakhena/gogen2/temp/infrastructure/log"
	"github.com/mirzaakhena/gogen2/temp/infrastructure/util"
	"github.com/mirzaakhena/gogen2/temp/usecase/createjournal"
	"net/http"
)

// createJournalHandler ...
func (r *Controller) createJournalHandler(inputPort createjournal.Inport) echo.HandlerFunc {

	return func(c echo.Context) error {

		//
		//
		//
		//
		//
		//
		//

		ctx := log.ContextWithLogGroupID(c.Request().Context())

		var req createjournal.InportRequest
		if err := c.Bind(&req); err != nil {
			newErr := apperror.FailUnmarshalResponseBodyError
			log.ErrorResponse(ctx, err)
			return c.JSON(http.StatusBadRequest, util.MustJSON(NewErrorResponse(newErr)))
		}

		log.InfoRequest(ctx, util.MustJSON(req))

		res, err := inputPort.Execute(ctx, req)
		if err != nil {
			log.ErrorResponse(ctx, err)
			return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
		}

		log.InfoResponse(ctx, util.MustJSON(res))
		return c.JSON(http.StatusOK, NewSuccessResponse(res))
	}
}
