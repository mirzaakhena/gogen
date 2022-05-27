package generror

import (
	"fmt"
	"github.com/mirzaakhena/gogen/utils"
	"go/token"
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

	domainName := utils.GetDefaultDomain()
	errorName := inputs[0]

	obj := ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		ErrorName:   errorName,
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	err = utils.CreateEverythingExactly("templates/", "errorenum", fileRenamer, struct{}{}, utils.AppTemplates)
	if err != nil {
		return err
	}

	errEnumFile := fmt.Sprintf("domain_%s/model/errorenum/error_enum.go", domainName)

	// inject to error_enum.go
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
