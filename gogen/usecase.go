package gogen

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
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

	{
		_, err := os.Stat(".application_schema/")
		if os.IsNotExist(err) {
			return fmt.Errorf("please call `gogen init` first")
		}
	}

	if len(args) < 3 {
		return fmt.Errorf("please define usecase name. ex: `gogen usecase CreateOrder`")
	}
	usecaseName := args[USECASE_NAME_INDEX]

	// if schema file is not found
	if _, err := os.Stat(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName)); os.IsNotExist(err) {

		WriteFile(
			".application_schema/usecases/usecase._yml",
			fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName),
			struct{}{},
		)

		CreateFolder("usecases/%s/inport", strings.ToLower(usecaseName))

		CreateFolder("usecases/%s/interactor", strings.ToLower(usecaseName))

		CreateFolder("usecases/%s/outport", strings.ToLower(usecaseName))

	}

	tp := readYAML(usecaseName)

	WriteFile(
		"usecases/usecase/inport/inport._go",
		fmt.Sprintf("usecases/%s/inport/inport.go", usecaseName),
		tp,
	)

	WriteFile(
		"usecases/usecase/outport/outport._go",
		fmt.Sprintf("usecases/%s/outport/outport.go", usecaseName),
		tp,
	)

	// check interactor file. only create if not exist
	if _, err := os.Stat(fmt.Sprintf("usecases/%s/interactor/interactor.go", usecaseName)); os.IsNotExist(err) {

		WriteFile(
			"usecases/usecase/interactor/interactor._go",
			fmt.Sprintf("usecases/%s/interactor/interactor.go", usecaseName),
			tp,
		)

	}

	// check interactor_test file. only create if not exist
	if _, err := os.Stat(fmt.Sprintf("usecases/%s/interactor/interactor_test.go", usecaseName)); os.IsNotExist(err) {
		WriteFile(
			"usecases/usecase/interactor/interactor_test._go",
			fmt.Sprintf("usecases/%s/interactor/interactor_test.go", usecaseName),
			tp,
		)
	}

	goFormat(tp.PackagePath)

	GenerateMock(tp.PackagePath, usecaseName)

	return nil
}

func readYAML(usecaseName string) *Usecase {

	content, err := ioutil.ReadFile(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName))
	if err != nil {
		log.Fatal(err)
	}

	tp := Usecase{}

	if err = yaml.Unmarshal(content, &tp); err != nil {
		log.Fatalf("error: %+v", err)
	}

	tp.Name = usecaseName
	tp.PackagePath = GetPackagePath()
	tp.Inport.RequestFieldObjs = ExtractField(tp.Inport.RequestFields)
	tp.Inport.ResponseFieldObjs = ExtractField(tp.Inport.ResponseFields)

	// x, _ := json.Marshal(tp.Inport)
	// fmt.Printf("%v\n", string(x))

	for i, out := range tp.Outports {
		tp.Outports[i].RequestFieldObjs = ExtractField(out.RequestFields)
		tp.Outports[i].ResponseFieldObjs = ExtractField(out.ResponseFields)

		// x, _ := json.Marshal(out)
		// fmt.Printf("%v\n", string(x))
	}

	return &tp

}

func ExtractField(fields []string) []Variable {

	vars := []Variable{}

	for _, field := range fields {
		s := strings.Split(field, " ")
		name := strings.TrimSpace(s[0])

		datatype := "string"
		if len(s) > 1 {
			datatype = strings.TrimSpace(s[1])
		}

		vars = append(vars, Variable{
			Name:     name,
			Datatype: datatype,
		})

	}

	return vars
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

func goFormat(path string) {
	fmt.Println("go fmt")
	cmd := exec.Command("go", "fmt", fmt.Sprintf("%s/...", path))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
