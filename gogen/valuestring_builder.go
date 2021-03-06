package gogen

import (
	"fmt"
	"strings"
)

type ValueStringBuilderRequest struct {
	ValueStringName string
	FolderPath      string
	GomodPath       string
}

type valueStringBuilder struct {
	ValueStringBuilderRequest ValueStringBuilderRequest
}

func NewValueString(req ValueStringBuilderRequest) Generator {
	return &valueStringBuilder{ValueStringBuilderRequest: req}
}

func (d *valueStringBuilder) Generate() error {

	valueStringName := strings.TrimSpace(d.ValueStringBuilderRequest.ValueStringName)
	folderPath := d.ValueStringBuilderRequest.FolderPath
	gomodPath := d.ValueStringBuilderRequest.GomodPath

	if len(valueStringName) == 0 {
		return fmt.Errorf("ValueStringName name must not empty")
	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	en := StructureValueString{
		ValueStringName: valueStringName,
		PackagePath:     packagePath,
	}

	CreateFolder("%s/domain/entity", folderPath)

	CreateFolder("%s/domain/repository", folderPath)

	CreateFolder("%s/domain/service", folderPath)

	_ = WriteFileIfNotExist(
		"domain/entity/valuestring._go",
		fmt.Sprintf("%s/domain/entity/%s.go", folderPath, PascalCase(valueStringName)),
		en,
	)

	_ = WriteFileIfNotExist(
		"domain/repository/repository._go",
		fmt.Sprintf("%s/domain/repository/repository._go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"domain/repository/database._go",
		fmt.Sprintf("%s/domain/repository/database._go", folderPath),
		struct{}{},
	)

	return nil
}
