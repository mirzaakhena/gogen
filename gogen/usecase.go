package gogen

import (
	"fmt"
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

		CreateFolder("usecases/%s/inport", uc.Name)
		WriteFile(
			"usecases/usecase/inport/inport._go",
			fmt.Sprintf("usecases/%s/inport/inport.go", uc.Name),
			uc,
		)

		CreateFolder("usecases/%s/outport", uc.Name)
		WriteFile(
			"usecases/usecase/outport/outport._go",
			fmt.Sprintf("usecases/%s/outport/outport.go", uc.Name),
			uc,
		)

		CreateFolder("usecases/%s/interactor", uc.Name)
		WriteFileIfNotExist(
			"usecases/usecase/interactor/interactor._go",
			fmt.Sprintf("usecases/%s/interactor/interactor.go", uc.Name),
			uc,
		)

		GenerateMock(packagePath, uc.Name)

		WriteFileIfNotExist(
			"usecases/usecase/interactor/interactor_test._go",
			fmt.Sprintf("usecases/%s/interactor/interactor_test.go", uc.Name),
			uc,
		)

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
