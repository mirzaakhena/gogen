package main

import (
	"fmt"
	"os"

	"github.com/mirzaakhena/gogen/gogen"
)

const message string = `try execute gogen by
  gogen <command>

some command available is
  init .
  model <model name>
	usecase <usecase name>
	outport usecase <method1> <method2> <method3> ...
  test <usecase name>
  datasource <datasource name> <usecase name>
  controller <controller type and framework> <usecase name>
  registry <controller type> <datasource name> <usecase name>

some controller type available is
  restapi.gin
  restapi.http

for some controller here is under development
  consumer.nsq
  grpc
	
usecase name is using pascal case
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

	case "test":

		gen := gogen.NewTest()
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

	case "registry":

		gen := gogen.NewRegistry()
		if err := gen.Generate(os.Args...); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(0)
		}

	case "model":

		gen := gogen.NewModel()
		if err := gen.Generate(os.Args...); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(0)
		}

	case "outport":

		gen := gogen.NewOutport()
		if err := gen.Generate(os.Args...); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(0)
		}

	default:
		fmt.Printf("command %s is not recognized. %s\n", command, message)
		os.Exit(0)

	}

}
