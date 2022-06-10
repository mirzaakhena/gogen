package application

import (
	"helloworld/domain_helloworld/controller/restapi"
	"helloworld/domain_helloworld/gateway/empty"
	"helloworld/domain_helloworld/usecase/gethello"
	"helloworld/shared/driver"
	"helloworld/shared/infrastructure/config"
	"helloworld/shared/infrastructure/logger"
	"helloworld/shared/infrastructure/server"
	"helloworld/shared/infrastructure/util"
)

type hello struct {
	httpHandler *server.GinHTTPHandler
	controller  driver.Controller
}

func (c hello) RunApplication() {
	c.controller.RegisterRouter()
	c.httpHandler.RunApplication()
}

func NewHello() func() driver.RegistryContract {

	const appName = "hello"

	return func() driver.RegistryContract {

		cfg := config.ReadConfig()

		appID := util.GenerateID(4)

		appData := driver.NewApplicationData(appName, appID)

		log := logger.NewSimpleJSONLogger(appData)

		httpHandler := server.NewGinHTTPHandler(log, cfg.Servers[appName].Address, appData)

		datasource := empty.NewGateway(log, appData, cfg)

		return &hello{
			httpHandler: &httpHandler,
			controller: &restapi.Controller{
				Log:            log,
				Config:         cfg,
				Router:         httpHandler.Router,
				GetHelloInport: gethello.NewUsecase(datasource),
			},
		}

	}
}
