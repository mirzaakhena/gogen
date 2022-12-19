package main

import (
	"flag"
	"fmt"
	"github.com/mirzaakhena/gogen/command/genapplication"
	"github.com/mirzaakhena/gogen/command/gencontroller"
	"github.com/mirzaakhena/gogen/command/gencrud"
	"github.com/mirzaakhena/gogen/command/gendomain"
	"github.com/mirzaakhena/gogen/command/genentity"
	"github.com/mirzaakhena/gogen/command/genenum"
	"github.com/mirzaakhena/gogen/command/generror"
	"github.com/mirzaakhena/gogen/command/gengateway"
	"github.com/mirzaakhena/gogen/command/genrepository"
	"github.com/mirzaakhena/gogen/command/genservice"
	"github.com/mirzaakhena/gogen/command/gentest"
	"github.com/mirzaakhena/gogen/command/genusecase"
	"github.com/mirzaakhena/gogen/command/genvalueobject"
	"github.com/mirzaakhena/gogen/command/genvaluestring"
)

var Version = "v0.0.1"

func main() {

	type C struct {
		Command string
		Func    func(...string) error
	}

	commands := make([]C, 0)

	commands = append(commands,
		C{"domain", gendomain.Run},
		C{"entity", genentity.Run},
		C{"valueobject", genvalueobject.Run},
		C{"valuestring", genvaluestring.Run},
		C{"enum", genenum.Run},
		C{"usecase", genusecase.Run},
		C{"repository", genrepository.Run},
		C{"service", genservice.Run},
		C{"test", gentest.Run},
		C{"gateway", gengateway.Run},
		C{"controller", gencontroller.Run},
		C{"error", generror.Run},
		C{"application", genapplication.Run},
		C{"crud", gencrud.Run},
	)

	commandMap := map[string]func(...string) error{}

	for _, c := range commands {
		commandMap[c.Command] = c.Func
	}

	//commandMap := map[string]func(...string) error{
	//	"domain":      gendomain.Run,      // dom
	//	"usecase":     genusecase.Run,     // uc
	//	"entity":      genentity.Run,      // ent
	//	"valueobject": genvalueobject.Run, // vo
	//	"valuestring": genvaluestring.Run, // vs
	//	"enum":        genenum.Run,        // enum
	//	"repository":  genrepository.Run,  // repo
	//	"service":     genservice.Run,     // svc
	//	"test":        gentest.Run,        // test
	//	"gateway":     gengateway.Run,     // gtw
	//	"controller":  gencontroller.Run,  // ctl
	//	"error":       generror.Run,       // err
	//	"application": genapplication.Run, // app
	//	"crud":        gencrud.Run,        // crud
	//	// "webapp":      genwebapp.Run,      //
	//	// "web":         genweb.Run,         // web
	//	// "openapi":     genopenapi.Run,     //
	//}

	flag.Parse()
	cmd := flag.Arg(0)

	if cmd == "" {
		fmt.Printf("Gogen %s\n", Version)
		fmt.Printf("Try one of this command to learn how to use it\n")
		for _, k := range commands {
			fmt.Printf("  gogen %s\n", k.Command)
		}
		return
	}

	var values = make([]string, 0)
	if flag.NArg() > 1 {
		values = flag.Args()[1:]
	}

	f, exists := commandMap[cmd]
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
