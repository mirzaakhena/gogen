package gogen

import (
	"fmt"
	"strings"
)

type ValueObjectBuilderRequest struct {
	ValueObjectName string
	FolderPath      string
	GomodPath       string
	FieldNames      []string
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
	fieldNames := d.ValueObjectBuilderRequest.FieldNames

	if len(valueObjectName) == 0 {
		return fmt.Errorf("ValueObject name name must not empty")
	}

	if len(fieldNames) < 1 {
		return fmt.Errorf("ValueObject at least have one field")
	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	en := StructureValueObject{
		PackagePath:     packagePath,
		ValueObjectName: valueObjectName,
		FieldNames:      fieldNames,
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
