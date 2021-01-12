package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mirzaakhena/gogen/gogen"
)

func usage() {

	const message string = `
try execute gogen by those command:



1.
gogen usecase <usecase name> [<method1> <method2> <method3> ... ]

This command is for creating new usecase
sample:

	gogen usecase CreateOrder
	or
	gogen usecase CreateOrder Publish Save

CreateOrder is a usecase name with PascalCase
Publish and Save is an Outport method

	
	
2.	
gogen outport usecase <method1> <method2> <method3> ...

This command is for adding new outport to existing usecase
sample:

	gogen outports CreateOrder Validate AnotherCheck
	
CreateOrder is a usecase name with PascalCase	
Validate and AnotherCheck is a new method in usecase's outport



3.
gogen test <usecase name>

This command is for create test template for specific usecase
sample:
		
	gogen test CreateOrder	

CreateOrder is a usecase name with PascalCase	



4.
gogen gateway <gateway name> <usecase name>

This command is for create gateway
sample:
	
	gogen gateway Production CreateOrder

Production is a gateway name
CreateOrder is a usecase name with PascalCase		



5.
gogen controller <controller name> <usecase name>

This command is for create controller
sample:

	gogen controller Restapi CreateOrder

Restapi is a controller name
CreateOrder is a usecase name with PascalCase			



6.
gogen registry <registry name> <controller name> <usecase name> <gateway name>

This command is for bind controller usecase and gateway
sample:

	gogen registry Default Restapi CreateOrder Production

Default is a registry name
Restapi is a controller name
CreateOrder is a usecase name with PascalCase
Production is a gateway name


`
	fmt.Fprintf(os.Stdout, "%s\n", message)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	var gen gogen.Generator

	switch flag.Arg(0) {

	case "usecase":

		if flag.NArg() == 1 {
			gogen.ShowAllUsecase()
			return
		}
		// gogen usecase CreateOrder Save Publish
		gen = gogen.NewUsecase(gogen.UsecaseBuilderRequest{
			FolderPath:         ".",
			UsecaseName:        flag.Arg(1),
			OutportMethodNames: flag.Args()[2:],
		})

	case "test":
		// gogen test CreateOrder
		gen = gogen.NewTest(gogen.TestBuilderRequest{
			FolderPath:  ".",
			UsecaseName: flag.Arg(1),
		})

	case "outports":
		// gogen outports CreateOrder Validate
		gen = gogen.NewOutport(gogen.OutportBuilderRequest{
			FolderPath:         ".",
			UsecaseName:        flag.Arg(1),
			OutportMethodNames: flag.Args()[2:],
		})

	case "gateway":
		// gogen gateway Production CreateOrder
		gen = gogen.NewGateway(gogen.GatewayBuilderRequest{
			FolderPath:  ".",
			GatewayName: flag.Arg(1),
			UsecaseName: flag.Arg(2),
		})

	case "controller":
		// gogen controller RestApi CreateOrder
		gen = gogen.NewController(gogen.ControllerBuilderRequest{
			FolderPath:     ".",
			ControllerName: flag.Arg(1),
			UsecaseName:    flag.Arg(2),
		})

	case "registry":
		// gogen registry Default RestApi CreateOrder Production
		gen = gogen.NewRegistry(gogen.RegistryBuilderRequest{
			FolderPath:     ".",
			RegistryName:   flag.Arg(1),
			ControllerName: flag.Arg(2),
			UsecaseName:    flag.Arg(3),
			GatewayName:    flag.Arg(4),
		})

	case "init":
		// gogen init
		gen = gogen.NewInit(gogen.InitBuilderRequest{
			FolderPath: ".",
		})

	default:
		usage()

	}

	if gen != nil {
		if err := gen.Generate(); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(0)
		}
	}

}
