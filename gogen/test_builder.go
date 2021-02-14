package gogen

import (
	"fmt"
	"strings"
)

type TestBuilderRequest struct {
	UsecaseName string
	FolderPath  string
	GomodPath   string
}

type testBuilder struct {
	TestBuilderRequest TestBuilderRequest
}

func NewTest(req TestBuilderRequest) Generator {
	return &testBuilder{TestBuilderRequest: req}
}

func (d *testBuilder) Generate() error {

	usecaseName := strings.TrimSpace(d.TestBuilderRequest.UsecaseName)
	folderPath := d.TestBuilderRequest.FolderPath
	gomodPath := d.TestBuilderRequest.GomodPath

	if len(usecaseName) == 0 {
		return fmt.Errorf("Usecase name must not empty")
	}

	// create a interactor_test.go file
	{
		mapStruct, errCollect := CollectPortStructs(folderPath, PascalCase(usecaseName))
		if errCollect != nil {
			return errCollect
		}

		uc, errConstruct := ConstructStructureUsecase(gomodPath, folderPath, PascalCase(usecaseName), mapStruct)
		if errConstruct != nil {
			return errConstruct
		}

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor_test._go",
			fmt.Sprintf("%s/usecase/%s/interactor_test.go", folderPath, LowerCase(usecaseName)),
			uc,
		)
	}

	GenerateMock(GetPackagePath(), PascalCase(usecaseName), folderPath)

	return nil
}
