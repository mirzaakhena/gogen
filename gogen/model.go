package gogen

type Variable struct {
	Name     string `yaml:"-"` //
	Datatype string `yaml:"-"` //
}

type Model struct {
	Name      string      `yaml:"name"`
	Fields    []string    `yaml:"fields"`
	FieldObjs []*Variable `yaml:"-"`
}

type Method struct {
	MethodName        string      `yaml:"methodName"`     //
	RequestFields     []string    `yaml:"requestFields"`  //
	ResponseFields    []string    `yaml:"responseFields"` //
	Models            []*Model    `yaml:"models"`         //
	RequestFieldObjs  []*Variable `yaml:"-"`              //
	ResponseFieldObjs []*Variable `yaml:"-"`              //
}

type Inport struct {
	RequestFields     []string    `yaml:"requestFields"`  //
	ResponseFields    []string    `yaml:"responseFields"` //
	Models            []*Model    `yaml:"models"`         //
	RequestFieldObjs  []*Variable `yaml:"-"`              //
	ResponseFieldObjs []*Variable `yaml:"-"`              //
}

type Outport struct {
	Name              string      `yaml:"name"`           //
	OutportExtends    []string    `yaml:"outportExtends"` // can inputed by Repository or Service interface
	RequestFields     []string    `yaml:"requestFields"`  //
	ResponseFields    []string    `yaml:"responseFields"` //
	Models            []*Model    `yaml:"models"`         //
	RequestFieldObjs  []*Variable `yaml:"-"`              //
	ResponseFieldObjs []*Variable `yaml:"-"`              //
}

type Usecase struct {
	Name           string     `yaml:"name"`     //
	Inport         *Inport    `yaml:"inport"`   //
	Outports       []*Outport `yaml:"outports"` //
	PackagePath    string     `yaml:"-"`        //
	DatasourceName string     `yaml:"-"`
}

type Service struct {
	Name           string    `yaml:"name"`           //
	ServiceMethods []*Method `yaml:"serviceMethods"` //
}

type Application struct {
	ApplicationName string        `yaml:"applicationName"` // name of application
	PackagePath     string        `yaml:"packagePath"`     // golang path of application
	Entities        []*Entity     `yaml:"entities"`        // list of entity used in this apps
	Usecases        []*Usecase    `yaml:"usecases"`        // list of usecase used in this apps
	Services        []*Service    `yaml:"services"`        // list of service used in this apps
	Repositories    []*Repository `yaml:"repositories"`    // list of repo used in this apps
}

type Entity struct {
	Name      string      `yaml:"name"`   // MANDATORY. name of the entity
	Fields    []string    `yaml:"fields"` // MANDATORY. all field under the entity
	FieldObjs []*Variable `yaml:"-"`      //
}

type Repository struct {
	Name        string   `yaml:"name"`        // MANDATORY. name of the Repo
	EntityName  string   `yaml:"entityName"`  //
	MethodNames []string `yaml:"methodNames"` //
	PackagePath string   `yaml:"-"`           //
}
