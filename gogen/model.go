package gogen

import (
	"fmt"
)

type model struct {
}

func NewModel() Generator {
	return &model{}
}

func (d *model) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define model name. ex: `gogen model Menu`")
	}

	return GenerateModel(ModelRequest{
		ModelName:  args[2],
		FolderPath: ".",
	})
}

type ModelRequest struct {
	ModelName  string
	FolderPath string
}

func GenerateModel(req ModelRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	CreateFolder("model/")

	_ = WriteFileIfNotExist(
		"model/model._go",
		fmt.Sprintf("%s/model/%s.go", req.FolderPath, req.ModelName),
		struct {
			Name      string
			Directory string
		}{
			Name:      req.ModelName,
			Directory: folderImport,
		},
	)

	return nil
}
