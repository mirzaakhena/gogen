package gogen

import (
	"fmt"
	"strings"
)

type EnumBuilderRequest struct {
	EnumName   string
	FolderPath string
}

type enumBuilder struct {
	EnumBuilderRequest EnumBuilderRequest
}

func NewEnum(req EnumBuilderRequest) Generator {
	return &enumBuilder{EnumBuilderRequest: req}
}

func (d *enumBuilder) Generate() error {

	enumName := strings.TrimSpace(d.EnumBuilderRequest.EnumName)
	folderPath := d.EnumBuilderRequest.FolderPath

	if len(enumName) == 0 {
		return fmt.Errorf("EnumName name must not empty")
	}

	en := StructureEnum{
		EnumName: enumName,
	}

	CreateFolder("%s/domain/entity", folderPath)

	CreateFolder("%s/domain/service", folderPath)

	_ = WriteFileIfNotExist(
		"domain/entity/enum._go",
		fmt.Sprintf("%s/domain/entity/%s.go", folderPath, PascalCase(enumName)),
		en,
	)

	return nil
}
