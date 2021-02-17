package gogen

import (
	"fmt"
	"strings"
)

type EnumBuilderRequest struct {
	FolderPath string
	EnumName   string
	EnumValues []string
	GomodPath  string
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
	enumValues := d.EnumBuilderRequest.EnumValues
	gomodPath := d.EnumBuilderRequest.GomodPath

	if len(enumName) == 0 {
		return fmt.Errorf("EnumName name must not empty")
	}

	if len(enumValues) < 2 {
		return fmt.Errorf("Enum at least have 2 value")
	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	en := StructureEnum{
		PackagePath: packagePath,
		EnumName:    enumName,
		EnumValues:  enumValues,
	}

	CreateFolder("%s/entity", folderPath)

	_ = WriteFileIfNotExist(
		"entity/enum._go",
		fmt.Sprintf("%s/entity/%s.go", folderPath, PascalCase(enumName)),
		en,
	)

	CreateFolder("%s/shared/errcat", folderPath)

	_ = WriteFileIfNotExist(
		"shared/errcat/error._go",
		fmt.Sprintf("%s/shared/errcat/error.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"shared/errcat/error_enum._go",
		fmt.Sprintf("%s/shared/errcat/error_enum.go", folderPath),
		struct{}{},
	)

	GoFormat(packagePath)

	return nil
}
