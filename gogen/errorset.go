package gogen

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type errorSet struct {
}

func NewErrorSet() Generator {
	return &errorSet{}
}

func (d *errorSet) Generate(args ...string) error {

	content, err := ioutil.ReadFile("errorset/error.yml")
	if err != nil {
		log.Fatal(err)
	}

	// prepare root object
	tp := ErrorSet{}

	// read yaml file
	if err = yaml.Unmarshal(content, &tp); err != nil {
		log.Fatalf("error: %+v", err)
	}

	_ = WriteFile("errorset/literal._go", "errorset/literal.go", &tp)

	_ = WriteFile("errorset/message._go", "errorset/message.go", &tp)

	_ = WriteFileIfNotExist("errorset/function._go", "errorset/function.go", struct{}{})

	GoFormat(GetPackagePath())

	return nil
}
