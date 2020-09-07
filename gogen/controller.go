package gogen

import (
	"fmt"
)

const (
	CONTROLLER_TYPE_INDEX         int = 2
	CONTROLLER_USECASE_NAME_INDEX int = 3
	// CONTROLLER_LIBRARY_INDEX      int = 3
	// CONTROLLER_USECASE_NAME_INDEX int = 4
)

type controller struct {
}

func NewController() Generator {
	return &controller{}
}

func (d *controller) Generate(args ...string) error {

	if IsNotExist(".application_schema/") {
		return fmt.Errorf("please call `gogen init` first")
	}

	if len(args) < 4 {
		return fmt.Errorf("please define controller_type, library/framework and usecase_name. ex: `gogen controller restapi CreateOrder`")
	}

	cType := args[CONTROLLER_TYPE_INDEX]

	// cLib := args[CONTROLLER_LIBRARY_INDEX]

	usecaseName := args[CONTROLLER_USECASE_NAME_INDEX]

	if IsNotExist(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName)) {
		return fmt.Errorf("Usecase `%s` is not found. Generate it by call `gogen usecase %s` first", usecaseName, usecaseName)
	}

	tp, err := ReadYAML(usecaseName)
	if err != nil {
		return err
	}

	if cType == "restapi" {

		CreateFolder("controllers/restapi")

		WriteFileIfNotExist(
			"controllers/restapi/gin._go",
			fmt.Sprintf("controllers/restapi/%s.go", usecaseName),
			tp,
		)

	} else if cType == "consumer" {

		CreateFolder("controllers/restapi")

		WriteFileIfNotExist(
			"controllers/restapi/gin._go",
			fmt.Sprintf("controllers/restapi/%s.go", usecaseName),
			tp,
		)

	}

	// if cType == "restapi" {

	// CreateFolder("controllers/restapi")

	// if cLib == "gin" {

	// WriteFileIfNotExist(
	// 	"controllers/restapi/gin._go",
	// 	fmt.Sprintf("controllers/restapi/%s.go", usecaseName),
	// 	tp,
	// )

	// 	} else //

	// 	if cLib == "echo" {

	// 		WriteFileIfNotExist(
	// 			"controllers/restapi/echo._go",
	// 			fmt.Sprintf("controllers/restapi/%s.go", usecaseName),
	// 			tp,
	// 		)

	// 	} else //

	// 	if cLib == "net/http" {

	// 		WriteFileIfNotExist(
	// 			"controllers/restapi/net-http._go",
	// 			fmt.Sprintf("controllers/restapi/%s.go", usecaseName),
	// 			tp,
	// 		)

	// 	} else //

	// 	if cLib == "mux" {

	// 		WriteFileIfNotExist(
	// 			"controllers/restapi/mux._go",
	// 			fmt.Sprintf("controllers/restapi/%s.go", usecaseName),
	// 			tp,
	// 		)

	// 	} else //

	// 	{
	// 		return fmt.Errorf("%s is not recognized. Only have gin, echo, net/http, mux", cLib)
	// 	}

	// } else //

	// if cType == "consumer" {

	// 	CreateFolder("controllers/consumer")

	// 	if cLib == "nsq" {

	// 		WriteFileIfNotExist(
	// 			"controllers/consumer/nsq._go",
	// 			fmt.Sprintf("controllers/consumer/%s.go", usecaseName),
	// 			tp,
	// 		)

	// 	} else //

	// 	if cLib == "rabbitmq" {

	// 		WriteFileIfNotExist(
	// 			"controllers/consumer/nsq._go",
	// 			fmt.Sprintf("controllers/consumer/%s.go", usecaseName),
	// 			tp,
	// 		)

	// 	} else //

	// 	if cLib == "kafka" {

	// 		WriteFileIfNotExist(
	// 			"controllers/consumer/kafka._go",
	// 			fmt.Sprintf("controllers/consumer/%s.go", usecaseName),
	// 			tp,
	// 		)

	// 	} else //

	// 	{
	// 		return fmt.Errorf("%s is not recognized. Only have nsq, rabbitmq, kafka", cLib)
	// 	}

	// } else //

	// if cType == "grpc" {

	// 	CreateFolder("controllers/grpc")

	// 	if cLib == "grpc" {
	// 		WriteFileIfNotExist(
	// 			"controllers/grpc/grpc._go",
	// 			fmt.Sprintf("controllers/grpc/%s.go", usecaseName),
	// 			tp,
	// 		)
	// 	} else //

	// 	if cLib == "net/rpc" {
	// 		WriteFileIfNotExist(
	// 			"controllers/grpc/rpc._go",
	// 			fmt.Sprintf("controllers/grpc/%s.go", usecaseName),
	// 			tp,
	// 		)
	// 	} else //

	// 	{
	// 		return fmt.Errorf("%s is not recognized. Only have grpc, net/rpc", cLib)
	// 	}

	// } else //

	// {
	// 	return fmt.Errorf("%s is not recognized. Only have restapi, consumer, grpc", cType)
	// }

	GoFormat(tp.PackagePath)

	return nil
}
