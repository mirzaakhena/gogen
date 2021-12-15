package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// authorized is an interceptor
func (r *Controller) authorized() gin.HandlerFunc {

	return func(c *gin.Context) {

		authorized := true

		if !authorized {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
