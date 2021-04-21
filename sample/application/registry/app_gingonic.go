package registry

import (
	"accounting/application"
	"accounting/controller/gingonic"
	"accounting/gateway"
	"accounting/infrastructure/server"
	"accounting/usecase/createjournal"
)

type appGinGonic struct {
	server.GinHTTPHandler
	gingonicController gingonic.Controller
	// TODO Another controller will added here ... <<<<<<
}

func NewAppGinGonic() application.RegistryContract {

	httpHandler := server.NewGinHTTPHandler("")
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

func (r *appGinGonic) SetupController() {
	r.gingonicController.RegisterRouter()
	// TODO another router call will added here ... <<<<<<
}
