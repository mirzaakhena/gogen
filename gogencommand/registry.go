package gogencommand

import (
	"flag"

	"github.com/mirzaakhena/gogen2/util"
)

type RegistryModel struct {
	RegistryName string
}

func NewRegistryModel() (Commander, error) {
	flag.Parse()

	//entityName := flag.Arg(1)
	//if entityName == "" {
	//	return nil, fmt.Errorf("entity name must not empty. `gogen entity RegistryName`")
	//}

	return &RegistryModel{
		//RegistryName: entityName,
	}, nil
}

func (obj *RegistryModel) Run() error {

	err := util.CreateFolderIfNotExist("application/registry")
	if err != nil {
		return err
	}

	//{
	//	outputFile := fmt.Sprintf("domain/entity/%s.go", strings.ToLower(obj.RegistryName))
	//	err = util.WriteFileIfNotExist(templates.RegistryFile, outputFile, obj)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil

}
