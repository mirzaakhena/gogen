package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	gogen entity Order
	gogen valueobject Name FirstName LastName
	gogen valuestring OrderID
	gogen enum PaymentMethod DANA Gopay Ovo
	gogen state OrderStatus WaitingPayment Complete Cancelled Expired
	gogen init
`
	fmt.Fprintf(os.Stdout, "%s\n", message)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	var gen gogen.Generator

	currentGopath := gogen.GetGopath()
	currentPath, _ := filepath.Abs("./")
	isUnderGopath := strings.HasPrefix(currentPath, currentGopath)

	folderPath := "."

	gomodPath := ""

	if !isUnderGopath {
		if !gogen.IsExist("go.mod") {
			fmt.Printf("%s\n", "go.mod file is not found. Please run 'go mod init your/go/project/' first")
			return
		}

		file, err := os.Open("go.mod")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			if strings.HasPrefix(row, "module") {
				moduleRow := strings.Split(row, " ")
				if len(moduleRow) > 1 {
					gomodPath = moduleRow[1]
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}

	switch flag.Arg(0) {

	case "usecase":
		// gogen usecase CreateOrder Save Publish
		gen = gogen.NewUsecase(gogen.UsecaseBuilderRequest{
			FolderPath:         folderPath,
			GomodPath:          gomodPath,
			UsecaseName:        flag.Arg(1),
			OutportMethodNames: flag.Args()[2:],
		})

	case "test":
		// gogen test CreateOrder
		gen = gogen.NewTest(gogen.TestBuilderRequest{
			FolderPath:  folderPath,
			GomodPath:   gomodPath,
			UsecaseName: flag.Arg(1),
		})

	case "outports":
		// gogen outports CreateOrder Validate
		gen = gogen.NewOutport(gogen.OutportBuilderRequest{
			FolderPath:         folderPath,
			GomodPath:          gomodPath,
			UsecaseName:        flag.Arg(1),
			OutportMethodNames: flag.Args()[2:],
		})

	case "gateway":
		// gogen gateway Production CreateOrder
		gen = gogen.NewGateway(gogen.GatewayBuilderRequest{
			FolderPath:  folderPath,
			GomodPath:   gomodPath,
			GatewayName: flag.Arg(1),
			UsecaseName: flag.Arg(2),
		})

	case "controller":
		// gogen controller RestApi CreateOrder
		gen = gogen.NewController(gogen.ControllerBuilderRequest{
			FolderPath:     folderPath,
			GomodPath:      gomodPath,
			ControllerName: flag.Arg(1),
			UsecaseName:    flag.Arg(2),
		})

	case "registry":
		// gogen registry Default RestApi CreateOrder Production
		gen = gogen.NewRegistry(gogen.RegistryBuilderRequest{
			FolderPath:     folderPath,
			GomodPath:      gomodPath,
			RegistryName:   flag.Arg(1),
			ControllerName: flag.Arg(2),
			UsecaseName:    flag.Arg(3),
			GatewayName:    flag.Arg(4),
		})

	case "entity":
		// gogen entity Order
		gen = gogen.NewEntity(gogen.EntityBuilderRequest{
			FolderPath: folderPath,
			EntityName: flag.Arg(1),
			GomodPath:  gomodPath,
		})

	case "valueobject":
		// gogen valueobject Name FirstName LastName
		gen = gogen.NewValueObject(gogen.ValueObjectBuilderRequest{
			FolderPath:      folderPath,
			ValueObjectName: flag.Arg(1),
			GomodPath:       gomodPath,
			FieldNames:      flag.Args()[2:],
		})

	case "valuestring":
		// gogen valuestring OrderID
		gen = gogen.NewValueString(gogen.ValueStringBuilderRequest{
			FolderPath:      folderPath,
			ValueStringName: flag.Arg(1),
			GomodPath:       gomodPath,
		})

	case "enum":
		// gogen enum PaymentMethod
		gen = gogen.NewEnum(gogen.EnumBuilderRequest{
			FolderPath: folderPath,
			EnumName:   flag.Arg(1),
			EnumValues: flag.Args()[2:],
			GomodPath:  gomodPath,
		})

	case "state":
		// gogen state OrderStatus
		gen = gogen.NewState(gogen.StateBuilderRequest{
			FolderPath:  folderPath,
			StateName:   flag.Arg(1),
			StateValues: flag.Args()[2:],
			GomodPath:   gomodPath,
		})

	case "error":
		// gogen error SomethingGoesWrongError
		gen = gogen.NewError(gogen.ErrorBuilderRequest{
			FolderPath: folderPath,
			ErrorName:  flag.Arg(1),
			GomodPath:  gomodPath,
		})

	case "init":
		// gogen init
		gen = gogen.NewInit(gogen.InitBuilderRequest{
			FolderPath: folderPath,
			GomodPath:  gomodPath,
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
