package gogencommand

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mirzaakhena/gogen/templates"
	"github.com/mirzaakhena/gogen/util"
)

type MethodModel struct {
	EntityName string
	MethodName string
}

func NewMethodModel() (Commander, error) {
	flag.Parse()

	values := flag.Args()[1:]
	if len(values) != 2 {
		return nil, fmt.Errorf("method format must include `gogen method MethodName EntityName`")
	}

	return &MethodModel{
		MethodName: values[0],
		EntityName: values[1],
	}, nil
}

func (obj *MethodModel) Run() error {

	// TODO check existency entity

	constTemplateCode, err := util.PrintTemplate(templates.MethodFile, obj)
	if err != nil {
		return err
	}

	existingFile := fmt.Sprintf("domain/entity/%s.go", strings.ToLower(obj.EntityName))
	file, err := os.Open(existingFile)
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

	buffer.WriteString(constTemplateCode)
	buffer.WriteString("\n")

	if err := ioutil.WriteFile(existingFile, buffer.Bytes(), 0644); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	// TODO inject to interactor

	return nil

}
