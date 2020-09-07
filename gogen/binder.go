package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

const (
	BINDER_USECASE_NAME_INDEX int = 2
)

type binder struct {
}

func NewBinder() Generator {
	return &binder{}
}

func (d *binder) Generate(args ...string) error {

	if IsNotExist(".application_schema/") {
		return fmt.Errorf("please call `gogen init` first")
	}

	if len(args) < 3 {
		return fmt.Errorf("please define usecase_name. ex: `gogen bind CreateOrder`")
	}

	usecaseName := args[BINDER_USECASE_NAME_INDEX]

	if IsNotExist(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName)) {
		return fmt.Errorf("Usecase `%s` is not found. Generate it by call `gogen usecase %s` first", usecaseName, usecaseName)
	}

	tp, err := ReadYAML(usecaseName)
	if err != nil {
		return err
	}

	if !UsecaseIsExist(usecaseName, tp) {
		InjectCode(usecaseName)
	}

	GoFormat(tp.PackagePath)

	return nil
}

func UsecaseIsExist(usecaseName string, data interface{}) bool {

	filename := fmt.Sprintf("%s/src/%s/binder/wiring_component.go", GetGopath(), GetPackagePath())
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	var buffer bytes.Buffer
	scanner := bufio.NewScanner(file)

	marker := fmt.Sprintf("// GOGEN_MARKER_BINDER_RESTAPI_GIN %s", PascalCase(usecaseName))
	codeInjection := "// GOGEN_CODE_INJECTION_BINDER_RESTAPI_GIN"
	for scanner.Scan() {
		row := scanner.Text()
		trimRow := strings.TrimSpace(row)
		if strings.HasPrefix(trimRow, marker) {
			return true
		}

		if strings.HasPrefix(trimRow, codeInjection) {

			func() {
				file, err := os.Open(fmt.Sprintf("%s/src/github.com/mirzaakhena/gogen/injection/BINDER_RESTAPI_GIN.txt", GetGopath()))
				if err != nil {
					return
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)

				for scanner.Scan() {
					row := scanner.Text()
					buffer.WriteString(row)
					buffer.WriteString("\n")
				}
			}()

			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")

	}

	tpl := template.Must(template.New("something").Funcs(FuncMap).Parse(buffer.String()))

	fileOut, err := os.Create(filename)
	if err != nil {
		return false
	}

	{
		err := tpl.Execute(fileOut, data)
		if err != nil {
			return false
		}
	}

	return false
}

func InjectCode(usecaseName string) {
	fmt.Printf("We will inject %s\n", usecaseName)
}
