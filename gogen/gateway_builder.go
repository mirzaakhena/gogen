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
	CreateFolder("%s/gateway/%s/repoimpl", folderPath, LowerCase(gatewayName))

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/gateway._go",
		fmt.Sprintf("%s/gateway/%s/%s.go", folderPath, gatewayName, PascalCase(usecaseName)),
		ds,
	)

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/repoimpl/database._go",
		fmt.Sprintf("%s/gateway/%s/repoimpl/database.go", folderPath, gatewayName),
		ds,
	)

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/repoimpl/transaction._go",
		fmt.Sprintf("%s/gateway/%s/repoimpl/transaction.go", folderPath, gatewayName),
		ds,
	)

	_ = WriteFileIfNotExist(
		"gateway/gatewayName/repoimpl/repository._go",
		fmt.Sprintf("%s/gateway/%s/repoimpl/repository.go", folderPath, gatewayName),
		ds,
	)

	d.createLog(folderPath)

	d.createUtil(folderPath)

	GoModTidy()
	GoImport(fmt.Sprintf("%s/gateway/%s/%s.go", folderPath, gatewayName, PascalCase(usecaseName)))

	return nil
}

func (d *gatewayBuilder) createUtil(folderPath string) {
	CreateFolder("%s/infrastructure/util", folderPath)

	_ = WriteFileIfNotExist(
		"infrastructure/util/utils._go",
		fmt.Sprintf("%s/infrastructure/util/utils.go", folderPath),
		struct{}{},
	)
}

func (d *gatewayBuilder) createLog(folderPath string) {
	CreateFolder("%s/infrastructure/log", folderPath)

	_ = WriteFileIfNotExist(
		"infrastructure/log/contract._go",
		fmt.Sprintf("%s/infrastructure/log/contract.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"infrastructure/log/implementation._go",
		fmt.Sprintf("%s/infrastructure/log/implementation.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"infrastructure/log/public._go",
		fmt.Sprintf("%s/infrastructure/log/public.go", folderPath),
		struct{}{},
	)
}
