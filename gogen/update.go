package gogen

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type generate struct {
}

func NewGenerate() Generator {
	return &generate{}
}

func (d *generate) Generate(args ...string) error {

	if IsNotExist("gogen_schema.yml") {
		return fmt.Errorf("please call `gogen init .` first")
	}

	app := Application{}
	{
		content, err := ioutil.ReadFile("gogen_schema.yml")
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal(content, &app); err != nil {
			return err
		}
	}

	// USECASE
	{
		for _, uc := range app.Usecases {

			uc.PackagePath = GetPackagePath()

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

			GenerateMock(app.PackagePath, uc.Name)

			WriteFileIfNotExist(
				"usecases/usecase/interactor/interactor_test._go",
				fmt.Sprintf("usecases/%s/interactor/interactor_test.go", uc.Name),
				uc,
			)

		}
	}

	// SERVICE
	{
		CreateFolder("services/")
		for _, sc := range app.Services {

			WriteFile(
				"services/service._go",
				fmt.Sprintf("services/%s.go", strings.ToLower(sc.Name)),
				sc,
			)

		}
	}

	// ENTITY
	{
		CreateFolder("entities/")
		for _, et := range app.Entities {

			et.FieldObjs = ExtractField(et.Fields)

			WriteFile(
				"entities/entity._go",
				fmt.Sprintf("entities/%s.go", strings.ToLower(et.Name)),
				et,
			)

		}

		for _, rp := range app.Repositories {

			rp.PackagePath = GetPackagePath()

			WriteFile(
				"repositories/repository._go",
				fmt.Sprintf("repositories/%s.go", strings.ToLower(rp.Name)),
				rp,
			)

		}
	}

	GoFormat(app.PackagePath)

	return nil
}
