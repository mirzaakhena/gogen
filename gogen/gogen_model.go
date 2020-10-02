package gogen

type Usecase struct {
	Name        string `yaml:"name"` //
	PackagePath string `yaml:"-"`    //
}

type Datasource struct {
	DatasourceName string   `yaml:"name"` //
	UsecaseName    string   ``            //
	PackagePath    string   `yaml:"-"`    //
	OutportMethods []string ``            //
}

type Test struct {
	UsecaseName    string   ``         //
	PackagePath    string   `yaml:"-"` //
	OutportMethods []string ``         //
}

type Controller struct {
	PackagePath  string      `yaml:"-"` //
	UsecaseName  string      ``         //
	InportFields []*NameType ``
}

type NameType struct {
	Name string
	Type string
}

type Model struct {
	Name string
}

// type Variable struct {
// 	Name     string `yaml:"-"` //
// 	Datatype string `yaml:"-"` //
// }

// type Model struct {
// 	Name string `yaml:"name"`
// }

// type Service struct {
// 	Name string `yaml:"name"` //
// }

// type Application struct {
// 	ApplicationName string        `yaml:"applicationName"` // name of application
// 	PackagePath     string        `yaml:"packagePath"`     // golang path of application
// 	Entities        []*Entity     `yaml:"entities"`        // list of entity used in this apps
// 	Usecases        []*Usecase    `yaml:"usecases"`        // list of usecase used in this apps
// 	Services        []*Service    `yaml:"services"`        // list of service used in this apps
// 	Repositories    []*Repository `yaml:"repositories"`    // list of repo used in this apps
// }

// type Entity struct {
// 	Name      string      `yaml:"name"`   // MANDATORY. name of the entity
// 	Fields    []string    `yaml:"fields"` // MANDATORY. all field under the entity
// 	FieldObjs []*Variable `yaml:"-"`      //
// }

// type Repository struct {
// 	Name        string   `yaml:"name"`        // MANDATORY. name of the Repo
// 	EntityName  string   `yaml:"entityName"`  //
// 	MethodNames []string `yaml:"methodNames"` //
// 	PackagePath string   `yaml:"-"`           //
// }
