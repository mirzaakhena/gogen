package gogen

import (
	"fmt"
	"strings"
)

type StateBuilderRequest struct {
	FolderPath  string
	StateName   string
	StateValues []string
	GomodPath   string
}

type stateBuilder struct {
	StateBuilderRequest StateBuilderRequest
}

func NewState(req StateBuilderRequest) Generator {
	return &stateBuilder{StateBuilderRequest: req}
}

func (d *stateBuilder) Generate() error {

	stateName := strings.TrimSpace(d.StateBuilderRequest.StateName)
	folderPath := d.StateBuilderRequest.FolderPath
	stateValues := d.StateBuilderRequest.StateValues
	gomodPath := d.StateBuilderRequest.GomodPath

	if len(stateName) == 0 {
		return fmt.Errorf("StateName name must not empty")
	}

	if len(stateValues) < 2 {
		return fmt.Errorf("State at least have 2 values")
	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	en := StructureState{
		StateName:   stateName,
		StateValues: stateValues,
		PackagePath: packagePath,
	}

	CreateFolder("%s/entity", folderPath)

	_ = WriteFileIfNotExist(
		"entity/state._go",
		fmt.Sprintf("%s/entity/%s.go", folderPath, PascalCase(stateName)),
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
