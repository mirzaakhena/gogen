package gogencommand

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mirzaakhena/gogen2/templates"
	"github.com/mirzaakhena/gogen2/util"
)

type EnumModel struct {
	PackagePath string
	EnumName    string
	EnumValues  []string
}

func NewEnumModel() (Commander, error) {
	flag.Parse()

	values := flag.Args()[1:]
	if len(values) < 2 {
		return nil, fmt.Errorf("enum at least have 2 values. Type with format : `gogen enum EnumName Enum1Value Enum2Value`")
	}

	return &EnumModel{
		PackagePath: util.GetGoMod(),
		EnumName:    values[0],
		EnumValues:  values[1:],
	}, nil
}

func (obj *EnumModel) Run() error {

	InitiateError()

	err := util.CreateFolderIfNotExist("domain/vo")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("domain/vo/%s.go", strings.ToLower(obj.EnumName))
		err = util.WriteFileIfNotExist(templates.EnumFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	util.GoFormat()

	return nil

}
