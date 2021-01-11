package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type GatewayBuilderRequest struct {
	UsecaseName string
	GatewayName string
	FolderPath  string
}

type gatewayBuilder struct {
	GatewayBuilderRequest GatewayBuilderRequest
}

func NewGateway(req GatewayBuilderRequest) Generator {
	return &gatewayBuilder{GatewayBuilderRequest: req}
}

func (d *gatewayBuilder) Generate() error {

	usecaseName := d.GatewayBuilderRequest.UsecaseName
	gatewayName := d.GatewayBuilderRequest.GatewayName
	folderPath := d.GatewayBuilderRequest.FolderPath

	if len(usecaseName) == 0 || len(gatewayName) == 0 {
		return fmt.Errorf("gogen gateway has 4 parameter. Try `gogen gateway prod yourUsecaseName`")
	}

	// create a gateway folder with gateway name
	CreateFolder("%s/gateway/%s", folderPath, strings.ToLower(gatewayName))

	var outportMethods []InterfaceMethod

	outportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName))
	fSet := token.NewFileSet()
	node, errParse := parser.ParseFile(fSet, outportFile, nil, parser.ParseComments)
	if errParse != nil {
		return errParse
	}

	mapStruct, errCollect := CollectPortStructs(folderPath, PascalCase(usecaseName))
	if errCollect != nil {
		return errCollect
	}

	var errRead error
	outportMethods, errRead = ReadInterfaceMethodAndField(node, fmt.Sprintf("%sOutport", PascalCase(usecaseName)), mapStruct)
	if errRead != nil {
		return errRead
	}

	ds := StructureGateway{
		GatewayName: gatewayName,
		PackagePath: GetPackagePath(),
		UsecaseName: usecaseName,
		Outport:     outportMethods,
	}

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/gateway._go",
		fmt.Sprintf("%s/gateway/%s/%s.go", folderPath, gatewayName, PascalCase(usecaseName)),
		ds,
	)

	return nil
}
