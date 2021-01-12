package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mirzaakhena/gogen/gogen"
)

func usage() {

	const message string = `try execute gogen by those command:
	gogen usecase CreateOrder
	gogen usecase CreateOrder Publish Save
	gogen outports CreateOrder Validate Publish Save
	gogen test CreateOrder
	gogen gateway production CreateOrder
	gogen controller restapi CreateOrder
	gogen registry default restapi CreateOrder production
	gogen init
`
	fmt.Fprintf(os.Stdout, "%s\n", message)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	var gen gogen.Generator

	switch flag.Arg(0) {

	case "usecase":
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
