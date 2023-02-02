package genwebapp

import (
	"fmt"

	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	EntityName string
	DomainName string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a complete CRUD webapp ui for specific entity\n" +
			"   gogen webapp Product\n" +
			"     'Product' is an existing entity name\n" +
			"\n")

		return err
	}

	gCfg := utils.GetGogenConfig()

	entityName := inputs[0]

	obj := &ObjTemplate{
		EntityName: entityName,
		DomainName: gCfg.Domain,
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(gCfg.Domain),
		"entityname": utils.SnakeCase(entityName),
	}

	err := utils.CreateEverythingExactly("templates/", "webapp", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	return nil

}
