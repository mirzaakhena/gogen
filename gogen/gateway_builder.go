package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type GatewayBuilderRequest struct {
	FolderPath  string
	GomodPath   string
	UsecaseName string
	GatewayName string
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
	gomodPath := d.GatewayBuilderRequest.GomodPath

	if len(usecaseName) == 0 || len(gatewayName) == 0 {
		return fmt.Errorf("gogen gateway has 4 parameter. Try `gogen gateway prod yourUsecaseName`")
	}

	var outportMethods []InterfaceMethod

	// to create a gateway, first we need to read the Outport
	node, err := parser.ParseFile(token.NewFileSet(), fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, LowerCase(usecaseName)), nil, parser.ParseComments)
	if err != nil {
		return err
	}

	mapStruct, err := CollectPortStructs(folderPath, PascalCase(usecaseName))
	if err != nil {
		return err
	}

	outportMethods, err = ReadInterfaceMethodAndField(node, fmt.Sprintf("%sOutport", PascalCase(usecaseName)), mapStruct)
	if err != nil {
		return err
	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	ds := StructureGateway{
		GatewayName: gatewayName,
		PackagePath: packagePath,
		UsecaseName: usecaseName,
		Outport:     outportMethods,
	}

	// create a gateway folder with gateway name
	CreateFolder("%s/gateway/%s/repositoryimpl", folderPath, LowerCase(gatewayName))

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/gateway._go",
		fmt.Sprintf("%s/gateway/%s/%s.go", folderPath, gatewayName, PascalCase(usecaseName)),
		ds,
	)

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/repositoryimpl/repository._go",
		fmt.Sprintf("%s/gateway/%s/repositoryimpl/repository.go", folderPath, gatewayName),
		ds,
	)

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/repositoryimpl/database._go",
		fmt.Sprintf("%s/gateway/%s/repositoryimpl/database.go", folderPath, gatewayName),
		ds,
	)

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/repositoryimpl/transaction._go",
		fmt.Sprintf("%s/gateway/%s/repositoryimpl/transaction.go", folderPath, gatewayName),
		ds,
	)

	createLog(folderPath)

	createUtil(folderPath)

	createDatabase(folderPath)

	//GoModTidy()
	GoImport(fmt.Sprintf("%s/gateway/%s/%s.go", folderPath, gatewayName, PascalCase(usecaseName)))

	return nil
}
