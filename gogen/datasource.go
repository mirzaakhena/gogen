package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type datasource struct {
}

func NewDatasource() Generator {
	return &datasource{}
}

func (d *datasource) Generate(args ...string) error {

	if len(args) < 4 {
		return fmt.Errorf("please define datasource and usecase_name. ex: `gogen datasource production CreateOrder`")
	}

	return GenerateDatasource(DatasourceRequest{
		DatasourceName: args[2],
		UsecaseName:    args[3],
		FolderPath:     ".",
	})

}

type DatasourceRequest struct {
	DatasourceName string
	UsecaseName    string
	FolderPath     string
}

func GenerateDatasource(req DatasourceRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	ds := Datasource{}
	ds.Directory = folderImport
	ds.DatasourceName = req.DatasourceName
	ds.UsecaseName = req.UsecaseName
	ds.PackagePath = GetPackagePath()

	{
		inportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName))
		node, errParse := parser.ParseFile(token.NewFileSet(), inportFile, nil, parser.ParseComments)
		if errParse != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}

		interfaceNames, err := ReadInterfaceMethodName(node, fmt.Sprintf("%s%s", req.UsecaseName, "Outport"))
		if err != nil {
			return fmt.Errorf("usecase %s is not found 111", req.UsecaseName)
		}

		for _, methodName := range interfaceNames {
			ds.Outports = append(ds.Outports, &Outport{
				Name: methodName,
			})

			for _, ot := range ds.Outports {
				ot.RequestFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", ot.Name, "Request"))
				ot.ResponseFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", ot.Name, "Response"))
			}

		}
	}

	CreateFolder("%s/datasource/%s", req.FolderPath, strings.ToLower(req.DatasourceName))

	_ = WriteFileIfNotExist(
		"datasource/datasourceName/datasource._go",
		fmt.Sprintf("%s/datasource/%s/%s.go", req.FolderPath, req.DatasourceName, req.UsecaseName),
		ds,
	)

	GoFormat(ds.PackagePath)

	return nil
}
