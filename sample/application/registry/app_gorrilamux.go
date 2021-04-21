package registry

import (
	"accounting/application"
	"accounting/controller/gorrilamux"
	"accounting/gateway"
	"accounting/infrastructure/server"
	"accounting/usecase/createjournal"
)

type appGorrilaMux struct {
	server.GorrilaMuxHandler
	gorrilamuxController gorrilamux.Controller
	// TODO Another controller will added here ... <<<<<<
}

func NewAppGorrilaMux() application.RegistryContract {

	httpHandler := server.NewGorrilaMuxHandler(":8080")
	datasource := gateway.NewInmemoryGateway()

	return &appGorrilaMux{
		GorrilaMuxHandler: httpHandler,
		gorrilamuxController: gorrilamux.Controller{
			Router:              httpHandler.Router.NewRoute(),
			CreateJournalInport: createjournal.NewUsecase(datasource),
			// TODO another Inport will added here ... <<<<<<
		},
		// TODO another controller will added here ... <<<<<<
	}
}

func (r *appGorrilaMux) SetupController() {
	r.gorrilamuxController.RegisterRouter()
	// TODO another router call will added here ... <<<<<<
}
