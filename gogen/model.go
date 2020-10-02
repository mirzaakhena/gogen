package gogen

import (
	"fmt"
	"strings"
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

	CreateFolder("model/")

	_ = WriteFileIfNotExist(
		"model/model._go",
		fmt.Sprintf("model/%s.go", strings.ToLower(modelName)),
		struct {
			Name string
		}{
			Name: modelName,
		},
	)

	return nil
}
