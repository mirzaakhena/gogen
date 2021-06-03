package registry

import (
	"accounting/application"
	"accounting/controller/gingonic"
	"accounting/gateway"
	"accounting/infrastructure/log"
	"accounting/infrastructure/loglib"
	"accounting/infrastructure/server"
	"accounting/usecase/createjournal"
)

type appGinGonic struct {
	server.GinHTTPHandler
	gingonicController gingonic.Controller
	// TODO Another controller will added here ... <<<<<<
}

func NewAppGinGonic() func() application.RegistryContract {
	return func() application.RegistryContract {

		log.SetLogPrinter(loglib.GetLogWithJSONFormat())

		httpHandler := server.NewGinHTTPHandler(":8080")
		datasource := gateway.NewInmemoryGateway()

		return &appGinGonic{
			GinHTTPHandler: httpHandler,
			gingonicController: gingonic.Controller{
				Router:              httpHandler.Router,
				CreateJournalInport: createjournal.NewUsecase(datasource),
				// TODO another Inport will added here ... <<<<<<
			},
			// TODO another controller will added here ... <<<<<<
		}

	}
}

func (r *appGinGonic) SetupController() {
	r.gingonicController.RegisterRouter()
	// TODO another router call will added here ... <<<<<<
}
