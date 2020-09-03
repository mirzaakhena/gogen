package gogen

import (
	"fmt"
)

const (
	USECASE_NAME_INDEX int = 2
)

type usecase struct {
}

func NewUsecase() Generator {
	return &usecase{}
}

func (d *usecase) Generate(args ...string) error {

	if IsNotExist(".application_schema/") {
		return fmt.Errorf("please call `gogen init` first")
	}

	if len(args) < 3 {
		return fmt.Errorf("please define usecase name. ex: `gogen usecase CreateOrder`")
	}
	usecaseName := args[USECASE_NAME_INDEX]

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

	GoFormat(tp.PackagePath)

	GenerateMock(tp.PackagePath, usecaseName)

	return nil
}

type Inport struct {
	RequestFields     []string   `yaml:"requestFields"`  //
	ResponseFields    []string   `yaml:"responseFields"` //
	RequestFieldObjs  []Variable ``                      //
	ResponseFieldObjs []Variable ``                      //
}

type Outport struct {
	Name              string     `yaml:"name"`           //
	OutportExtends    []string   `yaml:"outportExtends"` //
	RequestFields     []string   `yaml:"requestFields"`  //
	ResponseFields    []string   `yaml:"responseFields"` //
	RequestFieldObjs  []Variable ``                      //
	ResponseFieldObjs []Variable ``                      //
}

type Usecase struct {
	Name        string    ``
	PackagePath string    ``
	Inport      Inport    `yaml:"inport"`   //
	Outports    []Outport `yaml:"outports"` //
}

type Variable struct {
	Name     string
	Datatype string
}
