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
	// if IsNotExist("gogen_schema.yml") {
	// 	return fmt.Errorf("please call `gogen init .` first")
	// }

	if len(args) < 3 {
		return fmt.Errorf("please define usecase name. ex: `gogen usecase CreateOrder`")
	}

	usecaseName := args[2]

	// Read the structure
	// app := Application{}
	// {
	// 	app.PackagePath = GetPackagePath()
	// 	s := strings.Split(app.PackagePath, "/")
	// 	app.ApplicationName = s[len(s)-1]
	// }

	// app.Usecases = append(app.Usecases, &Usecase{
	// 	Name: usecaseName,
	// })

	// output, err := yaml.Marshal(app)
	// if err != nil {
	// 	return err
	// }

	// if err := ioutil.WriteFile("gogen_schema.yml", output, 0644); err != nil {
	// 	return err
	// }

	// USECASE
	{

		packagePath := GetPackagePath()

		uc := Usecase{
			Name:        usecaseName,
			PackagePath: packagePath,
		}

		CreateFolder("usecase/%s/port", strings.ToLower(uc.Name))

		WriteFile(
			"usecase/usecaseName/port/inport._go",
			fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(uc.Name)),
			uc,
		)

		WriteFile(
			"usecase/usecaseName/port/outport._go",
			fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(uc.Name)),
			uc,
		)

		WriteFileIfNotExist(
			"usecase/usecaseName/interactor._go",
			fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(uc.Name)),
			uc,
		)

		WriteFileIfNotExist(
			"usecase/usecaseName/interactor_test._go",
			fmt.Sprintf("usecase/%s/interactor_test.go", strings.ToLower(uc.Name)),
			uc,
		)

		GenerateMock(packagePath, uc.Name)
	}

	return nil
}

func Abc(args ...string) error {

	if IsNotExist(".application_schema/") {
		return fmt.Errorf("please call `gogen init` first")
	}

	if len(args) < 3 {
		return fmt.Errorf("please define usecase name. ex: `gogen usecase CreateOrder`")
	}
	usecaseName := args[2]

	// if schema file is not found
	WriteFileIfNotExist(
		".application_schema/usecases/usecase._yml",
		fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName),
		struct {
			Name string
		}{
			Name: usecaseName,
		},
	)

	tp, err := ReadYAML(usecaseName)
	if err != nil {
		return err
	}

	WriteFile(
		"inport/inport._go",
		fmt.Sprintf("inport/%s.go", usecaseName),
		tp,
	)

	WriteFile(
		"outport/outport._go",
		fmt.Sprintf("outport/%s.go", usecaseName),
		tp,
	)

	// check interactor file. only create if not exist
	WriteFileIfNotExist(
		"interactor/interactor._go",
		fmt.Sprintf("interactor/%s.go", usecaseName),
		tp,
	)

	// check interactor_test file. only create if not exist
	WriteFileIfNotExist(
		"interactor/interactor_test._go",
		fmt.Sprintf("interactor/%s_test.go", usecaseName),
		tp,
	)

	return nil
}
