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
		"usecase":     genusecase.Run,
		"entity":      genentity.Run,
		"valueobject": genvalueobject.Run,
		"valuestring": genvaluestring.Run,
		"enum":        genenum.Run,
		"repository":  genrepository.Run,
		"service":     genservice.Run,
		"gateway":     gengateway.Run,
		"controller":  gencontroller.Run,
		"error":       generror.Run,
		"test":        gentest.Run,
		"application": genapplication.Run,
		"crud":        gencrud.Run,
		"webapp":      genwebapp.Run,
		"web":         genweb.Run,
		"openapi":     genopenapi.Run,
		"domain":      gendomain.Run,
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
