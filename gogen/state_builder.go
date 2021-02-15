package gogen

import (
	"fmt"
	"strings"
)

type StateBuilderRequest struct {
	FolderPath  string
	StateName   string
	StateValues []string
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

	if len(stateName) == 0 {
		return fmt.Errorf("StateName name must not empty")
	}

	if len(stateValues) < 2 {
		return fmt.Errorf("State at least have 2 values")
	}

	en := StructureState{
		StateName:   stateName,
		StateValues: stateValues,
	}

	CreateFolder("%s/domain/state", folderPath)

	CreateFolder("%s/domain/service", folderPath)

	_ = WriteFileIfNotExist(
		"domain/state/state._go",
		fmt.Sprintf("%s/domain/state/%s.go", folderPath, PascalCase(stateName)),
		en,
	)

	return nil
}
