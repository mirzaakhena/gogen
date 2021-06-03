package gogencommand

import (
	"fmt"

	"github.com/mirzaakhena/gogen/templates"
	"github.com/mirzaakhena/gogen/util"
)

func InitiateHelper() error {

	err := util.CreateFolderIfNotExist("infrastructure/util")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("infrastructure/util/helper.go")
		err = util.WriteFileIfNotExist(templates.HelperFile, outputFile, struct{}{})
		if err != nil {
			return err
		}
	}

	return nil

}
