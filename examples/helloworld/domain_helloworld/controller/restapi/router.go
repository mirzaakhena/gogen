package restapi

import (
	"github.com/gin-gonic/gin"

	"helloworld/domain_helloworld/usecase/gethello"
	"helloworld/shared/infrastructure/config"
	"helloworld/shared/infrastructure/logger"
)

type Controller struct {
	Router         gin.IRouter
	Config         *config.Config
	Log            logger.Logger
	GetHelloInport gethello.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.Router.GET("/gethello", r.authorized(), r.getHelloHandler())
}
