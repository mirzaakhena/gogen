package generror

import (
	"fmt"
	"go/token"

	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	ErrorName   string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create an error enum\n" +
			"   gogen error SomethingGoesWrongError\n" +
			"     'SomethingGoesWrongError' is an error constant name\n" +
			"\n")

		return err
	}

	packagePath := utils.GetPackagePath()
	gcfg := utils.GetGogenConfig()
	errorName := inputs[0]

	obj := ObjTemplate{
		PackagePath: packagePath,
		ErrorName:   errorName,
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(gcfg.Domain),
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	err = utils.CreateEverythingExactly("templates/", "errorenum", fileRenamer, struct{}{}, utils.AppTemplates)
	if err != nil {
		return err
	}

	errEnumFile := fmt.Sprintf("domain_%s/model/errorenum/error_codes.go", gcfg)

	// inject to error_codes.go
	{
		fset := token.NewFileSet()
		utils.InjectToErrorEnum(fset, errEnumFile, errorName, "ER")
	}

	// reformat outport._go
	err = utils.Reformat(errEnumFile, nil)
	if err != nil {
		return err
	}

	return nil

}
