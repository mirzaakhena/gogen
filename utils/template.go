package utils

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed templates
var AppTemplates embed.FS

func PrintTemplate(templateString string, x interface{}) (string, error) {

	tpl, err := template.New("something").Funcs(FuncMap).Parse(templateString)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tpl.Execute(&buffer, x); err != nil {
		return "", err
	}

	return buffer.String(), nil

}
