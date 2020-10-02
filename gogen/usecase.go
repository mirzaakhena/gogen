package gogen

import (
	"fmt"
	"strings"
)

type usecase struct {
}

func NewUsecase() Generator {
	return &usecase{}
}

func (d *usecase) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define usecase name. ex: `gogen usecase CreateOrder`")
	}

	usecaseName := args[2]

	// USECASE
	{

		packagePath := GetPackagePath()

		uc := Usecase{
			Name:        usecaseName,
			PackagePath: packagePath,
		}

		CreateFolder("usecase/%s/port", strings.ToLower(uc.Name))

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/inport._go",
			fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(uc.Name)),
			uc,
		)

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/outport._go",
			fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(uc.Name)),
			uc,
		)

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor._go",
			fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(uc.Name)),
			uc,
		)

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor_test._go",
			fmt.Sprintf("usecase/%s/interactor_test.go", strings.ToLower(uc.Name)),
			uc,
		)

		GenerateMock(packagePath, uc.Name)
	}

	return nil
}
