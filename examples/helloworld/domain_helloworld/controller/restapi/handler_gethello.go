package restapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"helloworld/domain_helloworld/usecase/gethello"
	"helloworld/shared/infrastructure/logger"
	"helloworld/shared/infrastructure/util"
	"helloworld/shared/model/payload"
)

// getHelloHandler ...
func (r *Controller) getHelloHandler() gin.HandlerFunc {

	type request struct {
	}

	type response struct {
		Message string `json:"message"`
	}

	return func(c *gin.Context) {

		traceID := util.GenerateID(16)

		ctx := logger.SetTraceID(context.Background(), traceID)

		var jsonReq request
		if err := c.Bind(&jsonReq); err != nil {
			r.Log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		var req gethello.InportRequest

		r.Log.Info(ctx, util.MustJSON(req))

		res, err := r.GetHelloInport.Execute(ctx, req)
		if err != nil {
			r.Log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		var jsonRes response
		jsonRes.Message = res.Message

		r.Log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusOK, payload.NewSuccessResponse(jsonRes, traceID))

	}
}
