package gogencommand

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/mirzaakhena/gogen/templates"
	"github.com/mirzaakhena/gogen/util"
)

type TestModel struct {
	UsecaseName string
	TestName    string
	PackagePath string
	Methods     OutportMethods
}

func NewTestModel() (Commander, error) {
	flag.Parse()

	errMsg := fmt.Errorf("test format must include `gogen test testName UsecaseName`")

	// create a test need two params
	if flag.NArg() < 3 {
		return nil, errMsg
	}

	testName := flag.Arg(1)
	usecaseName := flag.Arg(2)

	return &TestModel{
		UsecaseName: usecaseName,
		TestName:    testName,
		PackagePath: util.GetGoMod(),
	}, nil
}

func (obj *TestModel) Run() error {

	err := obj.Methods.ReadOutport(obj.UsecaseName)
	if err != nil {
		return err
	}

	// Little hacky for replace the context.Context{} with ctx variable
	for _, m := range obj.Methods {
		if strings.HasPrefix(strings.TrimSpace(m.DefaultParamVal), `context.Context{}`) {
			m.DefaultParamVal = strings.ReplaceAll(m.DefaultParamVal, `context.Context{}`, "ctx")
		}
	}

	err = InitiateLog()
	if err != nil {
		return err
	}

	//err = obj.generateMock()
	//if err != nil {
	//	return err
	//}

	// create interactor_test.go
	{
		outputFile := fmt.Sprintf("usecase/%s/testcase_%s_test.go", strings.ToLower(obj.UsecaseName), strings.ToLower(obj.TestName))
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
