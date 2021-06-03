package gogencommand

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mirzaakhena/gogen2/templates"
	"github.com/mirzaakhena/gogen2/util"
)

type ValueStringModel struct {
	ValueStringName string
}

func NewValueStringModel() (Commander, error) {
	flag.Parse()

	voName := flag.Arg(1)
	if voName == "" {
		return nil, fmt.Errorf("valuestring name must not empty. `gogen valuestring VSName`")
	}

	return &ValueStringModel{
		ValueStringName: voName,
	}, nil
}

func (obj *ValueStringModel) Run() error {

	err := util.CreateFolderIfNotExist("domain/vo")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("domain/vo/%s.go", strings.ToLower(obj.ValueStringName))
		err = util.WriteFileIfNotExist(templates.ValueStringFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	return nil

}
