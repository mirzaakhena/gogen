package gogencommand

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mirzaakhena/gogen2/templates"
	"github.com/mirzaakhena/gogen2/util"
	"golang.org/x/tools/imports"
)

type UsecaseModel struct {
	UsecaseName string
	PackagePath string
}

func NewUsecaseModel() (Commander, error) {
	flag.Parse()

	// usecase only use one parameter which is the name of Usecase
	usecaseName := flag.Arg(1)
	if usecaseName == "" {
		return nil, fmt.Errorf("usecase name must not empty. `gogen usecase UsecaseName`")
	}

	return &UsecaseModel{
		UsecaseName: usecaseName,
		PackagePath: util.GetGoMod(),
	}, nil
}

func (obj *UsecaseModel) Run() error {

	// create a usecase and port folder
	err := util.CreateFolderIfNotExist("usecase/%s/port", strings.ToLower(obj.UsecaseName))
	if err != nil {
		return err
	}

	// create inport.go
	{
		outputFile := fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(obj.UsecaseName))
		err = util.WriteFileIfNotExist(templates.InportFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	// create outport.go
	{
		outputFile := fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(obj.UsecaseName))
		err = util.WriteFileIfNotExist(templates.OutportFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	// create interactor.go
	{
		outputFile := fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(obj.UsecaseName))
		err = util.WriteFileIfNotExist(templates.InteractorFile, outputFile, obj)
		if err != nil {
			return err
		}

		// reformat the import
		newBytes, err := imports.Process(outputFile, nil, nil)
		if err != nil {
			return err
		}

		// rewrite it
		if err := ioutil.WriteFile(outputFile, newBytes, 0644); err != nil {
			return err
		}
	}

	return nil

}
