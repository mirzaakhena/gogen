package gogen

type Usecase struct {
	Name                 string
	PackagePath          string
	Outports             []*Outport
	InportRequestFields  []*NameType
	InportResponseFields []*NameType
}

type Datasource struct {
	DatasourceName string
	UsecaseName    string
	PackagePath    string
	Outports       []*Outport
}

type Test struct {
	UsecaseName          string
	PackagePath          string
	Outports             []*Outport
	InportRequestFields  []*NameType
	InportResponseFields []*NameType
	Type                 string
}

type Controller struct {
	PackagePath string
	UsecaseName string
	// InportFields         []*NameType
	InportRequestFields  []*NameType
	InportResponseFields []*NameType
	Type                 string
}

type NameType struct {
	Name string
	Type string
}

type Model struct {
	Name string
}

type Outport struct {
	Name           string
	RequestFields  []*NameType
	ResponseFields []*NameType
}

type ErrorSet struct {
	Errors []ErrorItem `yaml:"errors"`
}

type ErrorItem struct {
	Code    int    `yaml:"code"`
	Literal string `yaml:"literal"`
	Message string `yaml:"message"`
}
