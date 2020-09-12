package gogen

import (
	"fmt"
	"io/ioutil"

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

			uc.Inport.RequestFieldObjs = ExtractField(uc.Inport.RequestFields)
			uc.Inport.ResponseFieldObjs = ExtractField(uc.Inport.ResponseFields)

			for i, out := range uc.Outports {
				uc.Outports[i].RequestFieldObjs = ExtractField(out.RequestFields)
				uc.Outports[i].ResponseFieldObjs = ExtractField(out.ResponseFields)
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
		for _, sc := range app.Services {

			serviceName := LowerCase(sc.Name)

			for _, sm := range sc.ServiceMethods {

				sm.RequestFieldObjs = ExtractField(sm.RequestFields)
				sm.ResponseFieldObjs = ExtractField(sm.ResponseFields)
			}

			CreateFolder("services/")
			WriteFile(
				"services/service._go",
				fmt.Sprintf("services/%s.go", serviceName),
				sc,
			)

		}
	}

	GoFormat(app.PackagePath)

	return nil
}
