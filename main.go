package main

import (
	"fmt"
	"os"

	"github.com/mirzaakhena/gogen/gogen"
)

const message string = `try execute gogen by
  gogen <command>

some command available is
  init
  usecase <usecase name>
  datasource <datasource name> <usecase name>
  controller <controller type> <usecase name>

some controller type available is
  restapi
  consumer
  grpc
`

func main() {

	arguments := os.Args
	if len(arguments) < 2 {

		fmt.Printf("%s\n", message)
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {

	case "init":

		gen := gogen.NewApplicationSchema()
		if err := gen.Generate(os.Args...); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(0)
		}

	case "usecase":

		gen := gogen.NewUsecase()
		if err := gen.Generate(os.Args...); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(0)
		}

	case "datasource":

		gen := gogen.NewDatasource()
		if err := gen.Generate(os.Args...); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(0)
		}

	case "controller":

		gen := gogen.NewController()
		if err := gen.Generate(os.Args...); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(0)
		}

	default:
		fmt.Printf("command %s is not recognized. %s\n", command, message)
		os.Exit(0)

	}

}
