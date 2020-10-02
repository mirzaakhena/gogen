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

	if len(args) < 4 {
		return fmt.Errorf("please define usecase name and type (command/query). ex: `gogen usecase command CreateOrder`")
	}

	usecaseType := args[2]

	usecaseName := args[3]

	// USECASE
	{

		packagePath := GetPackagePath()

		uc := Usecase{
			Name:        usecaseName,
			PackagePath: packagePath,
		}

		CreateFolder("usecase/%s/port", strings.ToLower(uc.Name))

		if usecaseType == "command" {

			_ = WriteFileIfNotExist(
				"usecase/usecaseName/port/inport-command._go",
				fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(uc.Name)),
				uc,
			)

			_ = WriteFileIfNotExist(
				"usecase/usecaseName/port/outport-command._go",
				fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(uc.Name)),
				uc,
			)

			_ = WriteFileIfNotExist(
				"usecase/usecaseName/interactor-command._go",
				fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(uc.Name)),
				uc,
			)

			// _ = WriteFileIfNotExist(
			// 	"usecase/usecaseName/interactor_test-command._go",
			// 	fmt.Sprintf("usecase/%s/interactor_test.go", strings.ToLower(uc.Name)),
			// 	uc,
			// )

		} else //

		if usecaseType == "query" {

			_ = WriteFileIfNotExist(
				"usecase/usecaseName/port/inport-query._go",
				fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(uc.Name)),
				uc,
			)

			_ = WriteFileIfNotExist(
				"usecase/usecaseName/port/outport-query._go",
				fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(uc.Name)),
				uc,
			)

			_ = WriteFileIfNotExist(
				"usecase/usecaseName/interactor-query._go",
				fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(uc.Name)),
				uc,
			)

			// _ = WriteFileIfNotExist(
			// 	"usecase/usecaseName/interactor_test-query._go",
			// 	fmt.Sprintf("usecase/%s/interactor_test.go", strings.ToLower(uc.Name)),
			// 	uc,
			// )

		} else //

		{
			return fmt.Errorf("use type `command` or `query`")
		}

		// GenerateMock(packagePath, uc.Name)
	}

	return nil
}
