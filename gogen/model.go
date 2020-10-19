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

	modelName := args[2]

	folderPath := "."

	return GenerateModel(modelName, folderPath)
}

func GenerateModel(modelName, folderPath string) error {

	var folderImport string
	if folderPath != "." {
		folderImport = fmt.Sprintf("/%s", folderPath)
	}

	CreateFolder("model/")

	_ = WriteFileIfNotExist(
		"model/model._go",
		fmt.Sprintf("%s/model/%s.go", folderPath, modelName),
		struct {
			Name      string
			Directory string
		}{
			Name:      modelName,
			Directory: folderImport,
		},
	)

	return nil
}
