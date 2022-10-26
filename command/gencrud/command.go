package gencrud

import (
	"fmt"
	"github.com/mirzaakhena/gogen/command/gengateway"
	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	EntityName  string
	DomainName  string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a new usecase\n" +
			"   gogen crud Product\n" +
			"     'Product' is an existing entity name\n" +
			"\n")

		return err
	}

	domainName := utils.GetDefaultDomain()

	entityName := inputs[0]

	obj := &ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		EntityName:  entityName,
		DomainName:  domainName,
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
		"entityname": utils.LowerCase(entityName),
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	// TODO check existing entity. if exist, read all the field else create new one

	err = utils.CreateEverythingExactly("templates/", "crud", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	err = gengateway.Run(fmt.Sprintf("gateway%s", entityName))
	if err != nil {
		return err
	}

	return nil

}
