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

	CreateFolder("%s/vo", folderPath)

	createDomain(folderPath)

	_ = WriteFileIfNotExist(
		"domain/vo/valuestring._go",
		fmt.Sprintf("%s/domain/vo/%s.go", folderPath, PascalCase(valueStringName)),
		en,
	)

	return nil
}
