package gogen

type Usecase struct {
	Name                 string     // name of usecase
	PackagePath          string     // root of your apps
	Directory            string     // apps directory after package path
	Outport              *Outport   //
	InportRequestFields  []NameType //
	InportResponseFields []NameType //
	DatasourceName       string     //
	CodeSlot             string     //
}

type Outport struct {
	UsecaseName string           //
	Methods     []*OutportMethod //
}

type OutportMethod struct {
	Name           string     // outport interface function's name
	RequestFields  []NameType //
	ResponseFields []NameType //
}

type NameType struct {
	Name    string //
	Type    string //
	Tag     string //
	Comment string //
}
