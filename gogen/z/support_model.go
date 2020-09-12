package gogen

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	TYPE_INDEX   int = 2
	TARGET_INDEX int = 3
	NAME_INDEX   int = 4
)

// gogen model inport  CreateOrder Something
// gogen model outport CreateOrder hello
// gogen model service CreateOrder Everything

type model struct {
}

func NewModel() Generator {
	return &model{}
}

func (d *model) Generate(args ...string) error {

	if IsNotExist("gogen_schema.yml") {
		return fmt.Errorf("please call `gogen init .` first")
	}

	typeModel := args[TYPE_INDEX]
	// targetModel := args[TARGET_INDEX]
	// nameModel := args[NAME_INDEX]

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

	if typeModel == "inport" {

		// for _, en := range app.SupportingModel.Inport {
		// 	if en.Name == targetModel {
		// 		return fmt.Errorf("Inport with name %s already exist", targetModel)
		// 	}
		// }

		// app.SupportingModel.Inport = append(app.SupportingModel.Inport, &Model{
		// 	Name:   nameModel,
		// 	Fields: []string{"field1 builtinType", "field2 OtherInportType"},
		// })

		output, err := yaml.Marshal(app)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile("gogen_schema.yml", output, 0644); err != nil {
			return err
		}

		return nil
	}

	if typeModel == "outport" {

		// for _, en := range app.SupportingModel.Outport {
		// 	if en.Name == targetModel {
		// 		return fmt.Errorf("Outport with name %s already exist", targetModel)
		// 	}
		// }

		// app.SupportingModel.Outport = append(app.SupportingModel.Outport, &Model{
		// 	Name:   nameModel,
		// 	Fields: []string{"field1 builtinType", "field2 OtherOutportType"},
		// })

		output, err := yaml.Marshal(app)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile("gogen_schema.yml", output, 0644); err != nil {
			return err
		}

		return nil
	}

	if typeModel == "service" {

		// for _, en := range app.SupportingModel.Service {
		// 	if en.Name == targetModel {
		// 		return fmt.Errorf("Service with name %s already exist", targetModel)
		// 	}
		// }

		// app.SupportingModel.Service = append(app.SupportingModel.Service, &Model{
		// 	Name:   nameModel,
		// 	Fields: []string{"field1 builtinType", "field2 OtherServiceType"},
		// })

		output, err := yaml.Marshal(app)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile("gogen_schema.yml", output, 0644); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("type must inport, outport or service")
}
