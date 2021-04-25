package gogencommand

import (
	"fmt"

	"github.com/mirzaakhena/gogen2/templates"
	"github.com/mirzaakhena/gogen2/util"
)

func InitiateLog() error {

	err := util.CreateFolderIfNotExist("infrastructure/log")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("infrastructure/log/log.go")
		err = util.WriteFileIfNotExist(templates.LogContractFile, outputFile, struct{}{})
		if err != nil {
			return err
		}
	}

	{
		outputFile := fmt.Sprintf("infrastructure/log/log_default.go")
		err = util.WriteFileIfNotExist(templates.LogImplFile, outputFile, struct{}{})
		if err != nil {
			return err
		}
	}

	return nil

}
