package gogen

type Usecase struct {
	Name                 string     // name of usecase
	PackagePath          string     // root of your apps
	Directory            string     // apps directory after package path
	Outports             []*Outport //
	InportRequestFields  []NameType //
	InportResponseFields []NameType //
}

type Outport struct {
	Name           string     // outport interface function's name
	RequestFields  []NameType //
	ResponseFields []NameType //
}

type Test struct {
	UsecaseName          string     // name of usecase
	PackagePath          string     // root of your apps
	Directory            string     // apps directory after package path
	Outports             []*Outport //
	InportRequestFields  []NameType //
	InportResponseFields []NameType //
}

type Controller struct {
	UsecaseName          string     // name of usecase
	PackagePath          string     // root of your apps
	Directory            string     // apps directory after package path
	InportRequestFields  []NameType //
	InportResponseFields []NameType //
}

type Datasource struct {
	DatasourceName string     //
	UsecaseName    string     //
	PackagePath    string     //
	Directory      string     //
	Outports       []*Outport //
}

type NameType struct {
	Name string //
	Type string //
}
