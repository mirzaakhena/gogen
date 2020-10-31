package gogen

import "fmt"

type schema struct {
}

func NewSchema() Generator {
	return &schema{}
}

func (d *schema) Generate(args ...string) error {

	if len(args) < 1 {
		return fmt.Errorf("")
	}

	return GenerateSchema(SchemaRequest{})
}

type SchemaRequest struct {
}

func GenerateSchema(req SchemaRequest) error {

	return nil
}
