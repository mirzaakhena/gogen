package registry

import (
	"accounting/application"
	"accounting/controller/julienschmidt"
	"accounting/gateway"
	"accounting/infrastructure/server"
	"accounting/usecase/createjournal"
)

type appJulienSchmidt struct {
	server.JulienSchmidtHandler
	julienschmidtController julienschmidt.Controller
	// TODO Another controller will added here ... <<<<<<
}

func NewAppJulienSchmidt() func() application.RegistryContract {

	return func() application.RegistryContract {

		httpHandler := server.NewJulienSchmidtHandler(":8080")
		datasource := gateway.NewInmemoryGateway()

		return &appJulienSchmidt{
			JulienSchmidtHandler: httpHandler,
			julienschmidtController: julienschmidt.Controller{
				Router:              httpHandler.Router,
				CreateJournalInport: createjournal.NewUsecase(datasource),
				// TODO another Inport will added here ... <<<<<<
			},
			// TODO another controller will added here ... <<<<<<
		}
	}
}

func (r *appJulienSchmidt) SetupController() {
	r.julienschmidtController.RegisterRouter()
	// TODO another router call will added here ... <<<<<<
}
