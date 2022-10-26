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
	"github.com/mirzaakhena/gogen/command/genopenapi"
	"github.com/mirzaakhena/gogen/command/genrepository"
	"github.com/mirzaakhena/gogen/command/genservice"
	"github.com/mirzaakhena/gogen/command/gentest"
	"github.com/mirzaakhena/gogen/command/genusecase"
	"github.com/mirzaakhena/gogen/command/genvalueobject"
	"github.com/mirzaakhena/gogen/command/genvaluestring"
	"github.com/mirzaakhena/gogen/command/genweb"
	"github.com/mirzaakhena/gogen/command/genwebapp"
)

func main() {

	commandMap := map[string]func(...string) error{
		"usecase":     genusecase.Run,     // uc
		"entity":      genentity.Run,      // ent
		"valueobject": genvalueobject.Run, // vo
		"valuestring": genvaluestring.Run, // vs
		"enum":        genenum.Run,        // enum
		"repository":  genrepository.Run,  // repo
		"service":     genservice.Run,     // svc
		"gateway":     gengateway.Run,     // gtw
		"controller":  gencontroller.Run,  // ctl
		"error":       generror.Run,       // err
		"test":        gentest.Run,        // test
		"application": genapplication.Run, // app
		"crud":        gencrud.Run,        // crud
		"webapp":      genwebapp.Run,      //
		"web":         genweb.Run,         // web
		"openapi":     genopenapi.Run,     //
		"domain":      gendomain.Run,      // dom
	}

	flag.Parse()
	cmd := flag.Arg(0)

	if cmd == "" {
		fmt.Printf("Try one of this command to learn how to use it\n")
		for k := range commandMap {
			fmt.Printf("  gogen %s\n", k)
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
