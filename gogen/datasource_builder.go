package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type DatasourceBuilderRequest struct {
	UsecaseName    string
	DatasourceName string
	FolderPath     string
}

type datasourceBuilder struct {
	DatasourceBuilderRequest DatasourceBuilderRequest
}

func NewDatasource(req DatasourceBuilderRequest) Generator {
	return &datasourceBuilder{DatasourceBuilderRequest: req}
}

func (d *datasourceBuilder) Generate() error {

	usecaseName := d.DatasourceBuilderRequest.UsecaseName
	datasourceName := d.DatasourceBuilderRequest.DatasourceName
	folderPath := d.DatasourceBuilderRequest.FolderPath

	// create a datasource folder with datasource name
	CreateFolder("%s/datasource/%s", folderPath, strings.ToLower(datasourceName))

	outportMethods := []InterfaceMethod{}
	{
		outportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName))
		fSet := token.NewFileSet()
		node, errParse := parser.ParseFile(fSet, outportFile, nil, parser.ParseComments)
		if errParse != nil {
			return errParse
		}

		mapStruct, errCollect := CollectPortStructs(folderPath, usecaseName)
		if errCollect != nil {
			return errCollect
		}

		var errRead error
		outportMethods, errRead = ReadInterfaceMethodAndField(node, fmt.Sprintf("%sOutport", usecaseName), mapStruct)
		if errRead != nil {
			return errRead
		}

	}

	ds := StructureDatasource{
		DatasourceName: datasourceName,
		PackagePath:    GetPackagePath(),
		UsecaseName:    usecaseName,
		Outport:        outportMethods,
	}

	_ = WriteFileIfNotExist(
		"datasource/datasourceName/datasource._go",
		fmt.Sprintf("%s/datasource/%s/%s.go", folderPath, datasourceName, usecaseName),
		ds,
	)

	return nil
}
