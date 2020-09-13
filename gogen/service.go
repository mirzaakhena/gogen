package gogen

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type service struct {
}

func NewService() Generator {
	return &service{}
}

func (d *service) Generate(args ...string) error {

	if IsNotExist("gogen_schema.yml") {
		return fmt.Errorf("please call `gogen init .` first")
	}

	serviceName := args[2]

	// Read gogen_schema and put it into object app
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

	for _, sc := range app.Services {
		if sc.Name == strings.TrimSpace(serviceName) {
			return fmt.Errorf("Service with name %s already exist", serviceName)
		}
	}

	var models []*Model
	models = append(models, &Model{
		Name:   "ModelName",
		Fields: []string{"Field1 BuiltinType", "Field2 OtherOutportType"},
	})

	var serviceMethods []*Method
	serviceMethods = append(serviceMethods, &Method{
		RequestFields:  []string{"Req1 string", "Req2 int"},
		ResponseFields: []string{"Res1 float64", "Res2 bool"},
		MethodName:     fmt.Sprintf("%sMethodName", serviceName),
		Models:         models,
	})

	app.Services = append(app.Services, &Service{
		Name:           serviceName,
		ServiceMethods: serviceMethods,
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
