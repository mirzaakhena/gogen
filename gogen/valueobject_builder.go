package gogen

import (
	"fmt"
	"strings"
)

type ValueObjectBuilderRequest struct {
	ValueObjectName string
	FolderPath      string
	GomodPath       string
}

type valueObjectBuilder struct {
	ValueObjectBuilderRequest ValueObjectBuilderRequest
}

func NewValueObject(req ValueObjectBuilderRequest) Generator {
	return &valueObjectBuilder{ValueObjectBuilderRequest: req}
}

func (d *valueObjectBuilder) Generate() error {

	valueObjectName := strings.TrimSpace(d.ValueObjectBuilderRequest.ValueObjectName)
	folderPath := d.ValueObjectBuilderRequest.FolderPath
	gomodPath := d.ValueObjectBuilderRequest.GomodPath

	if len(valueObjectName) == 0 {
		return fmt.Errorf("ValueObjectName name must not empty")
	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	en := StructureValueObject{
		ValueObjectName: valueObjectName,
		PackagePath:     packagePath,
	}

	CreateFolder("%s/domain/entity", folderPath)

	CreateFolder("%s/domain/service", folderPath)

	_ = WriteFileIfNotExist(
		"domain/entity/valueobject._go",
		fmt.Sprintf("%s/domain/entity/%s.go", folderPath, PascalCase(valueObjectName)),
		en,
	)

	return nil
}
