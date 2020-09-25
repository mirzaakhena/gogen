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
		return fmt.Errorf("please define datasource and usecase_name. ex: `gogen controller CreateOrder`")
	}

	usecaseName := args[2]

	ct := Controller{}
	ct.UsecaseName = usecaseName
	ct.PackagePath = GetPackagePath()

	CreateFolder("controllers/restapi")

	WriteFileIfNotExist(
		"controllers/restapi/gin._go",
		fmt.Sprintf("controllers/restapi/%s.go", usecaseName),
		ct,
	)

	GoFormat(ct.PackagePath)

	return nil
}
