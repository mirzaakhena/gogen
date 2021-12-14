package commandline

import (
	"github.com/mirzaakhena/gogen/usecase/gencontroller"
	"github.com/mirzaakhena/gogen/usecase/gencrud"
	"github.com/mirzaakhena/gogen/usecase/genentity"
	"github.com/mirzaakhena/gogen/usecase/generror"
	"github.com/mirzaakhena/gogen/usecase/gengateway"
	"github.com/mirzaakhena/gogen/usecase/genregistry"
	"github.com/mirzaakhena/gogen/usecase/genrepository"
	"github.com/mirzaakhena/gogen/usecase/genservice"
	"github.com/mirzaakhena/gogen/usecase/gentest"
	"github.com/mirzaakhena/gogen/usecase/genusecase"
	"github.com/mirzaakhena/gogen/usecase/genvalueobject"
	"github.com/mirzaakhena/gogen/usecase/genvaluestring"
	"github.com/mirzaakhena/gogen/usecase/genwebapp"
)

// Controller ...
type Controller struct {
	CommandMap           map[string]func(...string) error
	GenUsecaseInport     genusecase.Inport
	GenTestInport        gentest.Inport
	GenEntityInport      genentity.Inport
	GenRepositoryInport  genrepository.Inport
	GenServiceInport     genservice.Inport
	GenGatewayInport     gengateway.Inport
	GenErrorInport       generror.Inport
	GenControllerInport  gencontroller.Inport
	GenRegistryInport    genregistry.Inport
	GenValueObjectInport genvalueobject.Inport
	GenValueStringInport genvaluestring.Inport
	GenCrudInport        gencrud.Inport
	GenWebappInport      genwebapp.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	r.CommandMap["usecase"] = r.genUsecaseHandler(r.GenUsecaseInport)
	r.CommandMap["test"] = r.genTestHandler(r.GenTestInport)
	r.CommandMap["entity"] = r.genEntityHandler(r.GenEntityInport)
	r.CommandMap["repository"] = r.genRepositoryHandler(r.GenRepositoryInport)
	r.CommandMap["service"] = r.genServiceHandler(r.GenServiceInport)
	r.CommandMap["gateway"] = r.genGatewayHandler(r.GenGatewayInport)
	r.CommandMap["error"] = r.genErrorHandler(r.GenErrorInport)
	r.CommandMap["controller"] = r.genControllerHandler(r.GenControllerInport)
	r.CommandMap["registry"] = r.genRegistryHandler(r.GenRegistryInport)
	r.CommandMap["valueobject"] = r.genValueObjectHandler(r.GenValueObjectInport)
	r.CommandMap["valuestring"] = r.genValueStringHandler(r.GenValueStringInport)
	r.CommandMap["crud"] = r.genCrudHandler(r.GenCrudInport)
	r.CommandMap["webapp"] = r.genWebappHandler(r.GenWebappInport)
}
