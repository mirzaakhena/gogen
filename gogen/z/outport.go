package gogen

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type outport struct {
}

func NewOutport() Generator {
	return &outport{}
}

func (d *outport) Generate(args ...string) error {

	if IsNotExist("gogen_schema.yml") {
		return fmt.Errorf("please call `gogen init .` first")
	}

	usecaseName := args[2]
	outportName := args[3]

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

	var selectedUsecase *Usecase

	for _, uc := range app.Usecases {
		if uc.Name == usecaseName {
			selectedUsecase = uc
		}
	}
	if selectedUsecase == nil {
		return fmt.Errorf("Usecase with name %s is not found", usecaseName)
	}

	for _, op := range selectedUsecase.Outports {
		if op.Name == outportName {
			return fmt.Errorf("Outport with name %s already exist", outportName)
		}
	}

	output, err := yaml.Marshal(app)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("gogen_schema.yml", output, 0644); err != nil {
		return err
	}

	return nil
}
