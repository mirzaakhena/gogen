package registry

import (
	"flag"
	"fmt"
	"github.com/mirzaakhena/gogen/application"
	"github.com/mirzaakhena/gogen/controller"
	"github.com/mirzaakhena/gogen/controller/commandline"
	"github.com/mirzaakhena/gogen/gateway/prod"
	"github.com/mirzaakhena/gogen/usecase/gencontroller"
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
)

type gogenApp struct {
	CommandMap map[string]func(...string) error
	controller.Controller
}

// NewGogen2 ...
func NewGogen2() func() application.RegistryContract {
	return func() application.RegistryContract {

		datasource := prod.NewProdGateway()

		commandMap := make(map[string]func(...string) error, 0)

		return &gogenApp{
			CommandMap: commandMap,
			Controller: &commandline.Controller{
				CommandMap:           commandMap,
				GenUsecaseInport:     genusecase.NewUsecase(datasource),
				GenTestInport:        gentest.NewUsecase(datasource),
				GenEntityInport:      genentity.NewUsecase(datasource),
				GenRepositoryInport:  genrepository.NewUsecase(datasource),
				GenServiceInport:     genservice.NewUsecase(datasource),
				GenGatewayInport:     gengateway.NewUsecase(datasource),
				GenErrorInport:       generror.NewUsecase(datasource),
				GenControllerInport:  gencontroller.NewUsecase(datasource),
				GenRegistryInport:    genregistry.NewUsecase(datasource),
				GenValueObjectInport: genvalueobject.NewUsecase(datasource),
				GenValueStringInport: genvaluestring.NewUsecase(datasource),
			},
		}

	}
}

// RunApplication ...
func (r *gogenApp) RunApplication() {

	flag.Parse()
	cmd := flag.Arg(0)

	if cmd == "" {
		fmt.Printf("Those are the sample of gogen command: \n\n" +
			"   # Create a new usecase\n" +
			"   gogen usecase CreateOrder\n" +
			"     'CreateOrder' is an usecase name\n"+
			"\n"+
			"   # Create a test case file for current usecase\n" +
			"   gogen test normal CreateOrder\n" +
			"     'normal'      is a test case name\n"+
			"     'CreateOrder' is an usecase name\n"+
			"\n"+
			"   # Create a repository and inject the template code into interactor file with '//!' flag\n" +
			"   gogen repository SaveOrder Order CreateOrder\n" +
			"     'SaveOrder'   is a repository func name\n"+
			"     'Order'       is an entity name\n"+
			"     'CreateOrder' is an usecase name\n"+
			"\n"+
			"   # Create a repository without inject the template code into usecase\n" +
			"   gogen repository SaveOrder Order\n" +
			"     'SaveOrder' is a repository func name\n"+
			"     'Order'     is an entity name\n"+
			"\n"+
			"   # Create a service and inject the template code into interactor file with '//!' flag\n" +
			"   gogen service PublishMessage CreateOrder\n" +
			"     'PublishMessage' is a service func name\n"+
			"     'CreateOrder'    is an usecase name\n"+
			"\n"+
			"   # Create a service without inject the template code into usecase\n" +
			"   gogen service PublishMessage\n" +
			"     'PublishMessage' is a service func name\n"+
			"\n"+
			"   # Create a gateway for specific usecase\n" +
			"   gogen gateway inmemory CreateOrder\n" +
			"     'inmemory'    is a gateway name\n"+
			"     'CreateOrder' is an usecase name\n"+
			"\n"+
			"   # Create a gateway for all usecases\n" +
			"   gogen gateway inmemory\n" +
			"     'inmemory' is a gateway name\n"+
			"\n"+
			"   # Create a controller with defined web framework or other handler\n" +
			"   gogen controller restapi CreateOrder gin\n" +
			"     'restapi'     is a gateway name\n"+
			"     'CreateOrder' is an usecase name\n"+
			"     'gin'         is a sample webframewrok. You may try the other one like: nethttp, echo, and gorilla\n"+
			"\n"+
			"   # Create a controller with gin as default web framework\n" +
			"   gogen controller restapi CreateOrder\n" +
			"     'restapi'     is a gateway name\n"+
			"     'CreateOrder' is an usecase name\n"+
			"\n"+
			"   # Create a registry for specific controller\n" +
			"   gogen registry appone restapi\n" +
			"     'appone'  is an application name\n"+
			"     'restapi' is a controller name\n"+
			"\n"+
			"   # Create a registry for specific controller\n" +
			"   gogen entity Order\n" +
			"     'Order' is an entity name\n"+
			"\n"+
			"   # Create a valueobject with simple string type\n" +
			"   gogen valuestring OrderID\n" +
			"     'OrderID' is an valueobject name\n"+
			"\n"+
			"   # Create a valueobject with struct type\n" +
			"   gogen valueobject FullName FirstName LastName\n" +
			"     'FullName', 'FirstName', and 'LastName' is a Fields to created\n"+
			"\n"+
			"   # Create an enum predefine value\n" +
			"   gogen enum PaymentMethod DANA Gopay Ovo LinkAja\n" +
			"     'PaymentMethod' is an enum name\n"+
			"     'DANA', 'Gopay', 'Ovo', and 'LinkAja' is constant name\n"+
			"\n"+
			"   # Create an error enum\n" +
			"   gogen error SomethingGoesWrongError\n" +
			"     'SomethingGoesWrongError' is an error constant name\n"+
			"\n")
		return
	}

	var values = make([]string, 0)
	if flag.NArg() > 1 {
		values = flag.Args()[1:]
	}

	f, exists := r.CommandMap[cmd]
	if !exists {
		fmt.Printf("Command %s is not recognized\n", cmd)
		return
	}
	err := f(values...)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
}
