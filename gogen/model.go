package gogen

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
	Name           string    ``                //
	PackagePath    string    ``                //
	Inport         Inport    `yaml:"inport"`   //
	Outports       []Outport `yaml:"outports"` //
	DatasourceName string    ``
}

type Variable struct {
	Name     string `` //
	Datatype string `` //
}
