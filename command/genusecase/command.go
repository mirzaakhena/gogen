package genusecase

import (
	"fmt"
	"strings"

	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	UsecaseName string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a new usecase\n" +
			"   gogen usecase RunOrderCreate\n" +
			"     'RunOrderCreate' is an usecase name\n" +
			"\n")

		return err
	}

	packagePath := utils.GetPackagePath()
	domainName := utils.GetDefaultDomain()

	usecaseName := inputs[0]

	obj := &ObjTemplate{
		PackagePath: packagePath,
		UsecaseName: usecaseName,
	}

	fileRenamer := map[string]string{
		"usecasename": utils.LowerCase(usecaseName),
		"domainname":  domainName,
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	if strings.HasPrefix(utils.LowerCase(usecaseName), "getall") {
		err = utils.CreateEverythingExactly("templates/usecase/", "getall", fileRenamer, obj, utils.AppTemplates)
		if err != nil {
			return err
		}
	} else {
		err = utils.CreateEverythingExactly("templates/usecase/", "run", fileRenamer, obj, utils.AppTemplates)
		if err != nil {
			return err
		}
	}

	return nil

}
