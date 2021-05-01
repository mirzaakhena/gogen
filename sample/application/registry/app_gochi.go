package registry

import (
	"accounting/application"
	"accounting/controller/gochi"
	"accounting/gateway"
	"accounting/infrastructure/server"
	"accounting/usecase/createjournal"
)

type appGoChi struct {
	server.GoChiHandler
	gochiController gochi.Controller
	// TODO Another controller will added here ... <<<<<<
}

func NewAppGoChi() func() application.RegistryContract {

	return func() application.RegistryContract {

		httpHandler := server.NewGoChiHandler(":8080")
		datasource := gateway.NewInmemoryGateway()

		return &appGoChi{
			GoChiHandler: httpHandler,
			gochiController: gochi.Controller{
				Router:              httpHandler.Router,
				CreateJournalInport: createjournal.NewUsecase(datasource),
				// TODO another Inport will added here ... <<<<<<
			},
			// TODO another controller will added here ... <<<<<<
		}
	}
}

func (r *appGoChi) SetupController() {
	r.gochiController.RegisterRouter()
	// TODO another router call will added here ... <<<<<<
}
