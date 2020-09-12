package gogen

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type applicationSchema struct {
}

func NewApplicationSchema() Generator {
	return &applicationSchema{}
}

func (d *applicationSchema) Generate(args ...string) error {

	if IsExist("gogen_schema.yml") {

		app := Application{}
		{
			content, err := ioutil.ReadFile("gogen_schema.yml")
			if err != nil {
				return err
			}

			if err = yaml.Unmarshal(content, &app); err != nil {
				return fmt.Errorf("gogen_schema.yml exist but the format is invalid")
			}

		}

	} else {

		app := Application{}

		{
			app.PackagePath = GetPackagePath()
		}

		{
			s := strings.Split(app.PackagePath, "/")
			app.ApplicationName = s[len(s)-1]
		}

		output, err := yaml.Marshal(app)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile("gogen_schema.yml", output, 0644); err != nil {
			return err
		}

	}

	if len(args) < 3 {
		return fmt.Errorf("try `gogen init .` for current directory as active project\nor  `gogen init <your project name>` for create new project directory. But remember do `cd <your project name>` to start working")
	}

	directory := strings.TrimSpace(args[2])

	baseFolder := fmt.Sprintf("%s/", directory)

	folder := ""
	if directory != "." {
		folder = fmt.Sprintf("/%s", strings.TrimSpace(args[2]))
	}

	CreateFolder("%sbinder/", baseFolder)

	CreateFolder("%scontrollers/", baseFolder)

	CreateFolder("%sdatasources/", baseFolder)

	CreateFolder("%sentities/", baseFolder)

	CreateFolder("%susecases/", baseFolder)

	CreateFolder("%srepositories/", baseFolder)

	CreateFolder("%sservices/", baseFolder)

	CreateFolder("%sutils/", baseFolder)

	WriteFileIfNotExist(
		"main._go",
		fmt.Sprintf("%smain.go", baseFolder),
		struct {
			PackagePath string
			Directory   string
		}{
			PackagePath: GetPackagePath(),
			Directory:   folder,
		},
	)

	WriteFileIfNotExist(
		"config._toml",
		fmt.Sprintf("%sconfig.toml", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"README._md",
		fmt.Sprintf("%sREADME.md", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/runner._go",
		fmt.Sprintf("%sbinder/runner.go", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/setup_component._go",
		fmt.Sprintf("%sbinder/setup_component.go", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/wiring_component._go",
		fmt.Sprintf("%sbinder/wiring_component.go", baseFolder),
		struct {
			PackagePath string
		}{
			PackagePath: GetPackagePath(),
		},
	)

	return nil
}
