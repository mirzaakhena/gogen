package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type ControllerBuilderRequest struct {
	FolderPath     string
	GomodPath      string
	UsecaseName    string
	ControllerName string
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
	gomodPath := d.ControllerBuilderRequest.GomodPath

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

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	ct := StructureController{
		ControllerName: controllerName,
		PackagePath:    packagePath,
		UsecaseName:    usecaseName,
		Inport:         inportMethod,
	}

	CreateFolder("%s/infrastructure/log", folderPath)

	_ = WriteFileIfNotExist(
		"infrastructure/log/log._go",
		fmt.Sprintf("%s/infrastructure/log/log.go", folderPath),
		struct{}{},
	)

	CreateFolder("%s/infrastructure/util", folderPath)

	_ = WriteFileIfNotExist(
		"infrastructure/util/extractor_optimist._go",
		fmt.Sprintf("%s/infrastructure/util/extractor_optimist.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"infrastructure/util/extractor_pesimist._go",
		fmt.Sprintf("%s/infrastructure/util/extractor_pesimist.go", folderPath),
		struct{}{},
	)

	// create a controller folder with controller name
	CreateFolder("%s/controller/%s", folderPath, strings.ToLower(controllerName))

	_ = WriteFileIfNotExist(
		"controller/restapi/controller._go",
		fmt.Sprintf("%s/controller/%s/%s.go", folderPath, strings.ToLower(controllerName), PascalCase(usecaseName)),
		ct,
	)

	_ = WriteFileIfNotExist(
		"controller/interceptor._go",
		fmt.Sprintf("%s/controller/interceptor.go", folderPath),
		ct,
	)

	_ = WriteFileIfNotExist(
		"controller/response._go",
		fmt.Sprintf("%s/controller/response.go", folderPath),
		ct,
	)

	GoModTidy()

	return nil
}
