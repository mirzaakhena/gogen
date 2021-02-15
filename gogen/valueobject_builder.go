package gogen

import (
	"fmt"
	"strings"
)

type ValueObjectBuilderRequest struct {
	ValueObjectName string
	FolderPath      string
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

	if len(valueObjectName) == 0 {
		return fmt.Errorf("ValueObjectName name must not empty")
	}

	en := StructureValueObject{
		ValueObjectName: valueObjectName,
	}

	CreateFolder("%s/domain/valueObject", folderPath)

	CreateFolder("%s/domain/service", folderPath)

	_ = WriteFileIfNotExist(
		"domain/valueObject/valueObject._go",
		fmt.Sprintf("%s/domain/valueObject/%s.go", folderPath, PascalCase(valueObjectName)),
		en,
	)

	return nil
}
