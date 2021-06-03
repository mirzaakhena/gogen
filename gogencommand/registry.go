package gogencommand

import (
	"flag"
	"fmt"
	"github.com/mirzaakhena/gogen/templates"
	"github.com/mirzaakhena/gogen/util"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"strings"
)

type RegistryModel struct {
	PackagePath    string
	RegistryName   string
	ControllerName string
	UsecaseName    string
	GatewayName    string
}

func NewRegistryModel() (Commander, error) {
	flag.Parse()

	values := flag.Args()[1:]
	if len(values) != 4 {
		return nil, fmt.Errorf("entity name must not empty. `gogen registry RegistryName ControllerName UsecaseName GatewayName`")
	}

	return &RegistryModel{
		PackagePath:    util.GetGoMod(),
		RegistryName:   values[0],
		ControllerName: values[1],
		UsecaseName:    values[2],
		GatewayName:    values[3],
	}, nil
}

func (obj *RegistryModel) Run() error {

	err := util.CreateFolderIfNotExist("application/registry")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("application/application.go")
		err = util.WriteFileIfNotExist(templates.ApplicationFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	err = util.CreateFolderIfNotExist("infrastructure/server")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("infrastructure/server/http_handler.go")
		err = util.WriteFileIfNotExist(templates.ServerFile, outputFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(outputFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(outputFile, newBytes, 0644); err != nil {
			return err
		}
	}

	{
		outputFile := fmt.Sprintf("infrastructure/server/gracefully_shutdown.go")
		err = util.WriteFileIfNotExist(templates.ServerShutdownFile, outputFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(outputFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(outputFile, newBytes, 0644); err != nil {
			return err
		}
	}

	{
		outputFile := fmt.Sprintf("application/registry/%s.go", strings.ToLower(obj.RegistryName))
		err = util.WriteFileIfNotExist(templates.ApplicationRegistryGinFile, outputFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(outputFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(outputFile, newBytes, 0644); err != nil {
			return err
		}
	}

	{
		outputFile := fmt.Sprintf("main.go")
		err = util.WriteFileIfNotExist(templates.MainFile, outputFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(outputFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(outputFile, newBytes, 0644); err != nil {
			return err
		}
	}

	return nil

}
