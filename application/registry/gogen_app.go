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
		fmt.Printf("try \n\n" +
			"   gogen usecase CreateOrder\n" +
			"   gogen test normal CreateOrder\n" +
			"   gogen repository SaveOrder Order CreateOrder\n" +
			"   gogen repository SaveOrder Order\n" +
			"   gogen service PublishMessage CreateOrder\n" +
			"   gogen gateway inmemory CreateOrder\n" +
			"   gogen gateway inmemory\n" +
			"   gogen controller restapi CreateOrder gin\n" +
			"   gogen controller restapi CreateOrder\n" +
			"   gogen registry appone restapi\n" +
			"   gogen entity Order\n" +
			"   gogen valuestring OrderID\n" +
			"   gogen valueobject FullName FirstName LastName\n" +
			"   gogen enum PaymentMethod DANA Gopay Ovo LinkAja\n" +
			"   gogen error SomethingGoesWrongError\n" +
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
