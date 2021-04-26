package gogencommand

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/mirzaakhena/gogen2/templates"
	"github.com/mirzaakhena/gogen2/util"
)

type TestModel struct {
	UsecaseName string
	PackagePath string
	Methods     OutportMethods
}

func NewTestModel() (Commander, error) {
	flag.Parse()

	// create a test only need one params
	usecaseName := flag.Arg(1)
	if len(usecaseName) == 0 {
		return nil, fmt.Errorf("test format must include `gogen test UsecaseName`")
	}

	return &TestModel{
		UsecaseName: usecaseName,
		PackagePath: util.GetGoMod(),
	}, nil
}

func (obj *TestModel) Run() error {

	err := obj.Methods.ReadOutport(obj.UsecaseName)
	if err != nil {
		return err
	}

	// create interactor_test.go
	{
		outputFile := fmt.Sprintf("usecase/%s/interactor_test.go", strings.ToLower(obj.UsecaseName))
		if !util.IsExist(outputFile) {

			err = util.WriteFile(templates.TestFile, outputFile, obj)
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
	}

	err = obj.generateMock()
	if err != nil {
		return err
	}

	return nil
}

func (obj *TestModel) generateMock() error {

	usecaseName := strings.ToLower(obj.UsecaseName)

	cmd := exec.Command(
		"mockery",

		"--dir", fmt.Sprintf("usecase/%s/", usecaseName),

		// specifically interface with name that has suffix 'Outport'
		"--name", "Outport",

		// put the mock under usecase/%s/mocks/
		"-output", fmt.Sprintf("usecase/%s/mocks/", usecaseName),
	)

	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
