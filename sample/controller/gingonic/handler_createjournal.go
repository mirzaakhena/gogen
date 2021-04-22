package gingonic

import (
	"accounting/application/apperror"
	"accounting/infrastructure/log2"
	"accounting/infrastructure/util"
	"accounting/usecase/createjournal"
	"net/http"

	"github.com/gin-gonic/gin"
)

// createJournalHandler ...
func (r *Controller) createJournalHandler(inputPort createjournal.Inport) gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx := log2.Context(c.Request.Context())

		var req createjournal.InportRequest
		if err := c.BindJSON(&req); err != nil {
			newErr := apperror.FailUnmarshalResponseBodyError
			log2.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, NewErrorResponse(newErr))
			return
		}

		log2.Info(ctx, util.MustJSON(req))

		res, err := inputPort.Execute(ctx, req)
		if err != nil {
			log2.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, NewErrorResponse(err))
			return
		}

		log2.Info(ctx, util.MustJSON(res))
		c.JSON(http.StatusOK, NewSuccessResponse(res))

	}
}
