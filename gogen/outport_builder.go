package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

type OutportBuilderRequest struct {
	UsecaseName        string
	FolderPath         string
	GomodPath          string
	OutportMethodNames []string
}

type outportBuilder struct {
	OutportBuilderRequest OutportBuilderRequest
}

func NewOutport(req OutportBuilderRequest) Generator {
	return &outportBuilder{OutportBuilderRequest: req}
}

func (d *outportBuilder) Generate() error {

	usecaseName := d.OutportBuilderRequest.UsecaseName
	folderPath := d.OutportBuilderRequest.FolderPath
	newOutportMethodNames := d.OutportBuilderRequest.OutportMethodNames

	if len(usecaseName) == 0 {
		return fmt.Errorf("Usecase name must not empty")
	}

	if len(newOutportMethodNames) == 0 {
		return fmt.Errorf("Outport name is not defined. At least have one outport name")
	}

	// check the file outport.go exist or not
	outportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, LowerCase(usecaseName))
	if !IsExist(outportFile) {
		return fmt.Errorf("outport file is not found")
	}

	// prepare fileSet
	fSet := token.NewFileSet()
	node, errParse := parser.ParseFile(fSet, outportFile, nil, parser.ParseComments)
	if errParse != nil {
		return errParse
	}

	outportInterfaceName := fmt.Sprintf("%sOutport", PascalCase(usecaseName))
	outportInterfaceLine := fmt.Sprintf("type %s interface {", outportInterfaceName)

	oldMethodNames, _ := ReadInterfaceMethod(node, outportInterfaceName)

	file, err := os.Open(outportFile)
	if err != nil {
		return fmt.Errorf("not found outport file. You need to call 'gogen usecase %s' first", PascalCase(usecaseName))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	modeOutportInterface := false
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		if modeOutportInterface {

			if strings.HasPrefix(row, "}") {
				modeOutportInterface = false

				for _, methodName := range newOutportMethodNames {

					if _, exist := oldMethodNames[methodName]; exist {
						return fmt.Errorf("method with name %s already exist in interface %s", PascalCase(methodName), outportInterfaceName)
					}

					newMethodInterface, _ := PrintTemplate("usecase/usecaseName/port/outport_method._go", struct{ MethodName string }{MethodName: PascalCase(methodName)})
					buffer.WriteString(newMethodInterface)
					buffer.WriteString("\n")
				}

			}

		} else //

		if strings.HasPrefix(row, outportInterfaceLine) {
			modeOutportInterface = true
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")

	}

	for _, methodName := range newOutportMethodNames {
		newStruct, _ := PrintTemplate("usecase/usecaseName/port/outport_struct._go", struct{ MethodName string }{MethodName: PascalCase(methodName)})
		buffer.WriteString(newStruct)
		buffer.WriteString("\n")
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, LowerCase(usecaseName)), buffer.Bytes(), 0644); err != nil {
		return err
	}

	GoImport(fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, LowerCase(usecaseName)))

	return nil
}
