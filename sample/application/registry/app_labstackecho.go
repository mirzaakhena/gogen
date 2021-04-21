package registry

import (
	"accounting/application"
	"accounting/controller/labstackecho"
	"accounting/gateway"
	"accounting/infrastructure/server"
	"accounting/usecase/createjournal"
)

type appLabstackEcho struct {
	server.LabstackEchoHandler
	labstackechoController labstackecho.Controller
	// TODO Another controller will added here ... <<<<<<
}

func NewAppLabstackEcho() application.RegistryContract {

	httpHandler := server.NewLabstackEchoHandler(":8080")
	datasource := gateway.NewInmemoryGateway()

	return &appLabstackEcho{
		LabstackEchoHandler: httpHandler,
		labstackechoController: labstackecho.Controller{
			Router:              httpHandler.Router.Router(),
			CreateJournalInport: createjournal.NewUsecase(datasource),
			// TODO another Inport will added here ... <<<<<<
		},
		// TODO another controller will added here ... <<<<<<
	}
}

func (r *appLabstackEcho) SetupController() {
	r.labstackechoController.RegisterRouter()
	// TODO another router call will added here ... <<<<<<
}
