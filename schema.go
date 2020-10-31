package gogen

import "fmt"

type schema struct {
}

func NewSchema() Generator {
	return &schema{}
}

func (d *schema) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define model name. ex: `gogen model Menu`")
	}

	return GenerateSchema()
}

type SchemaRequest struct {
}

func GenerateSchema(req SchemaRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	CreateFolder("%s/model/", req.FolderPath)

	_ = WriteFileIfNotExist(
		"model/model._go",
		fmt.Sprintf("%s/model/%s.go", req.FolderPath, req.ModelName),
		struct {
			Name      string
			Directory string
		}{
			Name:      req.ModelName,
			Directory: folderImport,
		},
	)

	return nil
}
