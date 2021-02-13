package gogen

import (
	"fmt"
	"strings"
)

type EntityBuilderRequest struct {
	EntityName string
	FolderPath string
}

type entityBuilder struct {
	EntityBuilderRequest EntityBuilderRequest
}

func NewEntity(req EntityBuilderRequest) Generator {
	return &entityBuilder{EntityBuilderRequest: req}
}

func (d *entityBuilder) Generate() error {

	entityName := strings.TrimSpace(d.EntityBuilderRequest.EntityName)
	folderPath := d.EntityBuilderRequest.FolderPath

	if len(entityName) == 0 {
		return fmt.Errorf("EntityName name must not empty")
	}

	en := StructureEntity{
		EntityName: entityName,
	}

	_ = WriteFileIfNotExist(
		"domain/entity/entity._go",
		fmt.Sprintf("%s/domain/entity/%s.go", folderPath, strings.ToLower(entityName)),
		en,
	)

	return nil
}
