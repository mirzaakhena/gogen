package gogen

type Usecase struct {
	Name                 string
	PackagePath          string
	Directory            string
	Outports             []*Outport
	InportRequestFields  []NameType
	InportResponseFields []NameType
}

type Outport struct {
	Name           string
	RequestFields  []NameType
	ResponseFields []NameType
}

type Datasource struct {
	DatasourceName string
	UsecaseName    string
	PackagePath    string
	Directory      string
	Outports       []*Outport
}

type Test struct {
	UsecaseName          string
	PackagePath          string
	Directory            string
	Outports             []*Outport
	InportRequestFields  []NameType
	InportResponseFields []NameType
	Type                 string
}

type Controller struct {
	PackagePath          string
	UsecaseName          string
	Directory            string
	InportRequestFields  []NameType
	InportResponseFields []NameType
	Type                 string
}

type NameType struct {
	Name string
	Type string
}
