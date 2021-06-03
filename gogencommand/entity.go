package gogencommand

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mirzaakhena/gogen/templates"
	"github.com/mirzaakhena/gogen/util"
)

type EntityModel struct {
	EntityName string
}

func NewEntityModel() (Commander, error) {
	flag.Parse()

	entityName := flag.Arg(1)
	if entityName == "" {
		return nil, fmt.Errorf("entity name must not empty. `gogen entity EntityName`")
	}

	return &EntityModel{
		EntityName: entityName,
	}, nil
}

func (obj *EntityModel) Run() error {

	err := util.CreateFolderIfNotExist("domain/entity")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("domain/entity/%s.go", strings.ToLower(obj.EntityName))
		err = util.WriteFileIfNotExist(templates.EntityFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	return nil

}
