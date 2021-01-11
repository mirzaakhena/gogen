package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type ControllerBuilderRequest struct {
	UsecaseName    string
	ControllerName string
	FolderPath     string
	Framework      string
}

type controllerBuilder struct {
	ControllerBuilderRequest ControllerBuilderRequest
}

func NewController(req ControllerBuilderRequest) Generator {
	return &controllerBuilder{ControllerBuilderRequest: req}
}

func (d *controllerBuilder) Generate() error {

	usecaseName := d.ControllerBuilderRequest.UsecaseName
	controllerName := d.ControllerBuilderRequest.ControllerName
	folderPath := d.ControllerBuilderRequest.FolderPath
	framework := d.ControllerBuilderRequest.Framework

	if len(usecaseName) == 0 || len(controllerName) == 0 {
		return fmt.Errorf("gogen controller has 4 parameter. Try `gogen controller restapi yourUsecaseName`")
	}

	outportFile := fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, strings.ToLower(usecaseName))
	fSet := token.NewFileSet()
	node, errParse := parser.ParseFile(fSet, outportFile, nil, parser.ParseComments)
	if errParse != nil {
		return errParse
	}

	mapStruct, errCollect := CollectPortStructs(folderPath, PascalCase(usecaseName))
	if errCollect != nil {
		return errCollect
	}

	inportMethods, errRead := ReadInterfaceMethodAndField(node, fmt.Sprintf("%sInport", PascalCase(usecaseName)), mapStruct)
	if errRead != nil {
		return errRead
	}

	inportMethod := InterfaceMethod{}
	if len(inportMethods) == 1 {
		inportMethod = inportMethods[0]
	}

	ct := StructureController{
		ControllerName: controllerName,
		PackagePath:    GetPackagePath(),
		UsecaseName:    usecaseName,
		Inport:         inportMethod,
	}

	// create a controller folder with controller name
	CreateFolder("%s/controller/%s", folderPath, strings.ToLower(controllerName))

	if framework == "nethttp" {
		_ = WriteFileIfNotExist(
			"controller/restapi/controller_http._go",
			fmt.Sprintf("%s/controller/%s/%s.go", folderPath, strings.ToLower(controllerName), PascalCase(usecaseName)),
			ct,
		)

		_ = WriteFileIfNotExist(
			"controller/interceptor_http._go",
			fmt.Sprintf("%s/controller/interceptor.go", folderPath),
			ct,
		)
	} else //

	if framework == "gin" {
		_ = WriteFileIfNotExist(
			"controller/restapi/controller_gin._go",
			fmt.Sprintf("%s/controller/%s/%s.go", folderPath, strings.ToLower(controllerName), PascalCase(usecaseName)),
			ct,
		)

		_ = WriteFileIfNotExist(
			"controller/interceptor_gin._go",
			fmt.Sprintf("%s/controller/interceptor.go", folderPath),
			ct,
		)

	} else //

	{
		return fmt.Errorf("not recognize framework")
	}

	return nil
}
