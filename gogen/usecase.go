package gogen

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type usecase struct {
}

func NewUsecase() Generator {
	return &usecase{}
}

func (d *usecase) Generate(args ...string) error {
	if IsNotExist("gogen_schema.yml") {
		return fmt.Errorf("please call `gogen init .` first")
	}

	if len(args) < 3 {
		return fmt.Errorf("please define usecase name. ex: `gogen usecase CreateOrder`")
	}

	usecaseName := args[2]

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

	for _, uc := range app.Usecases {
		if uc.Name == strings.TrimSpace(usecaseName) {
			return fmt.Errorf("Usecase with name %s already exist", usecaseName)
		}
	}

	var inportModels []*Model
	inportModels = append(inportModels, &Model{
		Name:   "modelName",
		Fields: []string{"field1 builtinType", "field2 OtherOutportType"},
	})

	inport := Inport{
		RequestFields:  []string{"Req1 string", "Req2 int"},
		ResponseFields: []string{"Res1 float64", "Res2 bool"},
		Models:         inportModels,
	}

	var outportModels []*Model
	outportModels = append(outportModels, &Model{
		Name:   "modelName",
		Fields: []string{"field1 builtinType", "field2 OtherOutportType"},
	})

	var outports []*Outport
	outports = append(outports, &Outport{
		Name:           "OutportName",
		RequestFields:  []string{"Req1 string", "Req2 int"},
		ResponseFields: []string{"Res1 float64", "Res2 bool"},
		OutportExtends: []string{"Service1", "Repo1"},
		Models:         outportModels,
	})

	app.Usecases = append(app.Usecases, &Usecase{
		Name:     usecaseName,
		Inport:   &inport,
		Outports: outports,
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
