package gogencommand

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mirzaakhena/gogen2/templates"

	"github.com/mirzaakhena/gogen2/util"
	"golang.org/x/tools/imports"
)

type GatewayModel struct {
	PackagePath string         //
	UsecaseName string         //
	GatewayName string         //
	Methods     OutportMethods //
}

func NewGatewayModel() (Commander, error) {
	flag.Parse()

	// create a gateway need 2 parameter
	values := flag.Args()[1:]
	if len(values) != 2 {
		return nil, fmt.Errorf("gateway name must not empty. `gogen gateway GatewayName UsecaseName`")
	}

	return &GatewayModel{
		PackagePath: util.GetGoMod(),
		GatewayName: values[0],
		UsecaseName: values[1],
		Methods:     OutportMethods{},
	}, nil
}

func (obj *GatewayModel) Run() error {

	err := InitiateError()
	if err != nil {
		return err
	}

	err = InitiateLog()
	if err != nil {
		return err
	}

	err = InitiateHelper()
	if err != nil {
		return err
	}

	err = obj.Methods.ReadOutport(obj.UsecaseName)
	//err = obj.ReadOutport()
	if err != nil {
		return err
	}

	// create a gateway folder
	err = util.CreateFolderIfNotExist("gateway")
	if err != nil {
		return err
	}

	gatewayFile := fmt.Sprintf("gateway/%s.go", strings.ToLower(obj.GatewayName))

	// gateway file not exist yet
	if !util.IsExist(gatewayFile) {
		err := util.WriteFile(templates.GatewayFile, gatewayFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(gatewayFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(gatewayFile, newBytes, 0644); err != nil {
			return err
		}

	} else //

	// already exist. check existency current function implementation
	{

		existingFunc, err := (&GatewayReader{PackagePath: obj.PackagePath}).ReadCurrentGateway(obj.GatewayName)
		if err != nil {
			return err
		}

		// collect the only methods that has not added yet
		notExistingMethod := make([]*method, 0)
		for _, m := range obj.Methods {
			if _, exist := existingFunc[m.MethodName]; !exist {
				notExistingMethod = append(notExistingMethod, m)
			}
		}

		// we will only inject the non existing method
		obj.Methods = notExistingMethod
		{

			// reopen the file
			file, err := os.Open(gatewayFile)
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(file)
			var buffer bytes.Buffer
			for scanner.Scan() {
				row := scanner.Text()

				buffer.WriteString(row)
				buffer.WriteString("\n")
			}

			if err := file.Close(); err != nil {
				return err
			}

			constTemplateCode, err := util.PrintTemplate(templates.GatewayMethodFile, obj)
			if err != nil {
				return err
			}

			// write the template in the end of file
			buffer.WriteString(constTemplateCode)
			buffer.WriteString("\n")

			// reformat and do import
			newBytes, err := imports.Process(gatewayFile, buffer.Bytes(), nil)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(gatewayFile, newBytes, 0644); err != nil {
				return err
			}

		}

	}

	return nil

}
