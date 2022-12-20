package genenum

import (
	"fmt"

	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	EnumName    string
	EnumValues  []string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create an enum with values\n" +
			"   gogen enum PaymentMethod DANA Gopay Ovo\n" +
			"     'PaymentMethod'           is a enum name to created\n" +
			"     'DANA', 'Gopay' and 'Ovo' is an enum values\n" +
			"\n")
		return err
	}

	packagePath := utils.GetPackagePath()
	gfcg := utils.GetGogenConfig()
	enumName := inputs[0]

	obj := &ObjTemplate{
		PackagePath: packagePath,
		EnumName:    enumName,
		EnumValues:  inputs[1:],
	}

	fileRenamer := map[string]string{
		"enumname":   utils.SnakeCase(enumName),
		"domainname": utils.LowerCase(gfcg.Domain),
	}

	err := utils.CreateEverythingExactly("templates/", "enum", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	return nil

}
