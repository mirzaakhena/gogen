package gogen

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// content, err := ioutil.ReadFile(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName))
// if err != nil {
// 	log.Fatal(err)
// 	return nil, fmt.Errorf("cannot read %s.yml", usecaseName)
// }

// tp := Usecase{}

// if err = yaml.Unmarshal(content, &tp); err != nil {
// 	log.Fatalf("error: %+v", err)
// 	return nil, fmt.Errorf("%s.yml is unrecognized usecase file", usecaseName)
// }

func NewTest() Generator {
	return &mytest{}
}

type mytest struct {
}

func (d *mytest) Generate(args ...string) error {

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

	app.Entities = append(app.Entities, &Entity{
		Name: "Order",
	})

	app.ApplicationName = "Something"

	output, err := yaml.Marshal(app)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("gogen_schema.yml", output, 0644); err != nil {
		return err
	}

	return nil
}
