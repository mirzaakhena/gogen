package genusecase

import (
	"fmt"
	"github.com/mirzaakhena/gogen/utils"
	"strings"
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
			"   gogen usecase CreateOrder\n" +
			"     'CreateOrder' is an usecase name\n" +
			"\n")

		return err
	}

	domainName := utils.GetDefaultDomain()

	usecaseName := inputs[0]

	obj := &ObjTemplate{
		PackagePath: utils.GetPackagePath(),
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
