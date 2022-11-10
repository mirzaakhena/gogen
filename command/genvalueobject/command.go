package genvalueobject

import (
	"fmt"

	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath     string
	ValueObjectName string
	FieldNames      []string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a valueobject with struct type\n" +
			"   gogen valueobject FullName FirstName LastName\n" +
			"     'FullName'                 is a Value Object name to created\n" +
			"     'FirstName' and 'LastName' is a Fields on Value Object\n" +
			"\n")
		return err
	}

	packagePath := utils.GetPackagePath()
	domainName := utils.GetDefaultDomain()
	valueObjectName := inputs[0]

	obj := &ObjTemplate{
		PackagePath:     packagePath,
		ValueObjectName: valueObjectName,
		FieldNames:      inputs[1:],
	}

	fileRenamer := map[string]string{
		"valueobjectname": utils.SnakeCase(valueObjectName),
		"domainname":      utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "valueobject", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	return nil

}
