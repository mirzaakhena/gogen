package gogen

import (
	"fmt"
)

type controller struct {
}

func NewController() Generator {
	return &controller{}
}

func (d *controller) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define datasource and usecase_name. ex: `gogen controller restapi.gin CreateOrder`")
	}

	controllerType := args[2]

	usecaseName := args[3]

	ct := Controller{}
	ct.UsecaseName = usecaseName
	ct.PackagePath = GetPackagePath()

	if controllerType == "restapi.gin" {

		CreateFolder("controller/restapi")

		WriteFileIfNotExist(
			"controller/restapi/gin._go",
			fmt.Sprintf("controller/restapi/%s.go", usecaseName),
			ct,
		)

	} else //

	if controllerType == "restapi.http" {

		CreateFolder("controller/restapi")

		WriteFileIfNotExist(
			"controller/restapi/http._go",
			fmt.Sprintf("controller/restapi/%s.go", usecaseName),
			ct,
		)

	}

	GoFormat(ct.PackagePath)

	return nil
}
