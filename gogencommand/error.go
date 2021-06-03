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
)

type ErrorModel struct {
	ErrorName string
}

func NewErrorModel() (Commander, error) {
	flag.Parse()

	errorName := flag.Arg(1)
	if errorName == "" {
		return nil, fmt.Errorf("error name must not empty. `gogen error ErrorName`")
	}

	return &ErrorModel{
		ErrorName: errorName,
	}, nil
}

func (obj *ErrorModel) Run() error {

	err := InitiateError()
	if err != nil {
		return err
	}

	constTemplateCode, err := util.PrintTemplate(templates.ErrorEnumTemplateFile, obj)
	if err != nil {
		return err
	}

	existingFile := fmt.Sprintf("application/apperror/error_enum.go")
	file, err := os.Open(existingFile)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	var constFound bool = false
	for scanner.Scan() {
		row := scanner.Text()

		if constFound && strings.HasPrefix(row, ")") {
			constFound = false
			buffer.WriteString(constTemplateCode)
			buffer.WriteString("\n")

		} else //

		if strings.HasPrefix(row, "const (") {
			constFound = true

		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	if err := ioutil.WriteFile(existingFile, buffer.Bytes(), 0644); err != nil {
		return err
	}

	util.GoImport(existingFile)

	if err := file.Close(); err != nil {
		return err
	}

	return nil

}

func InitiateError() error {
	// create error

	err := util.CreateFolderIfNotExist("application/apperror")
	if err != nil {
		return err
	}

	{
		outputFile := fmt.Sprintf("application/apperror/error_enum.go")
		err = util.WriteFileIfNotExist(templates.ErrorEnumFile, outputFile, struct{}{})
		if err != nil {
			return err
		}
	}

	{
		outputFile := fmt.Sprintf("application/apperror/error_func.go")
		err = util.WriteFileIfNotExist(templates.ErrorFuncFile, outputFile, struct{}{})
		if err != nil {
			return err
		}
	}

	return nil

}
