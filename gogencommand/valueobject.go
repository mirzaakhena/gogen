package gogencommand

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mirzaakhena/gogen/templates"
	"github.com/mirzaakhena/gogen/util"
)

type ValueObjectModel struct {
	ValueObjectName string
	FieldNames      []string
}

func NewValueObjectModel() (Commander, error) {
	flag.Parse()

	values := flag.Args()[1:]
	if len(values) < 2 {
		return nil, fmt.Errorf("valueobject at least have 2 fields with format : `gogen valueobject VOName Field1 Field2`")
	}

	return &ValueObjectModel{
		ValueObjectName: values[0],
		FieldNames:      values[1:],
	}, nil
}

func (obj *ValueObjectModel) Run() error {

	err := util.CreateFolderIfNotExist("domain/vo")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("domain/vo/%s.go", strings.ToLower(obj.ValueObjectName))
		err = util.WriteFileIfNotExist(templates.ValueObjectFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	return nil

}
