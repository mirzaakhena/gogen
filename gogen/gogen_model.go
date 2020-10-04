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
	PackagePath  string
	UsecaseName  string
	InportFields []*NameType
	Type         string
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
