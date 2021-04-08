package gogencommand

import (
	"flag"
	"fmt"
	"github.com/mirzaakhena/gogen2/templates"
	"github.com/mirzaakhena/gogen2/util"
	"strings"
)

type TestModel struct {
	UsecaseName string
	PackagePath string
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

	// create a mock folder
	err := util.CreateFolderIfNotExist("usecase/mock")
	if err != nil {
		return err
	}

	// create interactor_test.go
	{
		outputFile := fmt.Sprintf("usecase/%s/interactor_test.go", strings.ToLower(obj.UsecaseName))
		err = util.WriteFileIfNotExist(templates.TestFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	return nil
}
