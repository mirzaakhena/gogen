package registry

import (
	"accounting/application"
	"accounting/controller/nethttp"
	"accounting/gateway"
	"accounting/infrastructure/server"
	"accounting/usecase/createjournal"
)

type appNetHttp struct {
	server.NetHTTPHandler
	nethttpController nethttp.Controller
	// TODO Another controller will added here ... <<<<<<
}

func NewAppNetHttp() application.RegistryContract {

	httpHandler := server.NewNetHTTPHandler("")
	datasource := gateway.NewInmemoryGateway()

	return &appNetHttp{
		NetHTTPHandler: httpHandler,
		nethttpController: nethttp.Controller{
			Router:              httpHandler.Router,
			CreateJournalInport: createjournal.NewUsecase(datasource),
			// TODO another Inport will added here ... <<<<<<
		},
		// TODO another controller will added here ... <<<<<<
	}
}

func (r *appNetHttp) SetupController() {
	r.nethttpController.RegisterRouter()
	// TODO another router call will added here ... <<<<<<
}
