package gogen

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type entity struct {
}

func NewEntity() Generator {
	return &entity{}
}

func (d *entity) Generate(args ...string) error {

	if IsNotExist("gogen_schema.yml") {
		return fmt.Errorf("please call `gogen init .` first")
	}

	if len(args) < 3 {
		return fmt.Errorf("please define entity name. ex: `gogen entity Menu`")
	}

	entityName := args[2]

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

	for _, et := range app.Entities {
		if et.Name == entityName {
			return fmt.Errorf("Entity with name %s already exist", entityName)
		}
	}

	app.Entities = append(app.Entities, &Entity{
		Name:   entityName,
		Fields: []string{"Field1 string", "Field2 int"},
	})

	app.Repositories = append(app.Repositories, &Repository{
		Name:        fmt.Sprintf("%sRepo", entityName),
		EntityName:  entityName,
		MethodNames: []string{fmt.Sprintf("%sRepoSave", entityName)},
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
