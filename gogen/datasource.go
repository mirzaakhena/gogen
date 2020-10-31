package gogen

import (
	"fmt"
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

	uc := Usecase{
		DatasourceName: req.DatasourceName,
		Name:           req.UsecaseName,
		Directory:      folderImport,
		PackagePath:    GetPackagePath(),
		Outport: &Outport{
			UsecaseName: req.UsecaseName,
		},
	}

	if err := readOutport(uc.Outport, req.FolderPath, req.UsecaseName); err != nil {
		return err
	}

	CreateFolder("%s/datasource/%s", req.FolderPath, strings.ToLower(req.DatasourceName))

	_ = WriteFileIfNotExist(
		"datasource/datasourceName/datasource._go",
		fmt.Sprintf("%s/datasource/%s/%s.go", req.FolderPath, req.DatasourceName, req.UsecaseName),
		uc,
	)

	GoFormat(uc.PackagePath)

	return nil
}
