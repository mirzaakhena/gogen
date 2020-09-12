package gogen

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type enum struct {
}

func NewEnum() Generator {
	return &enum{}
}

func (d *enum) Generate(args ...string) error {

	if IsNotExist("gogen_schema.yml") {
		return fmt.Errorf("please call `gogen init .` first")
	}

	enumName := args[2]

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

	for _, en := range app.Enums {
		if en.Name == enumName {
			return fmt.Errorf("Enum with name %s already exist", enumName)
		}
	}

	app.Enums = append(app.Enums, &Enum{
		Name:   enumName,
		Values: []string{"StateOne", "StateTwo"},
	})

	output, err := yaml.Marshal(app)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("gogen_schema.yml", output, 0644); err != nil {
		return err
	}

	return nil
}
