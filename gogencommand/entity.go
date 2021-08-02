package gogencommand

import (
	"flag"
	"fmt"
	"github.com/mirzaakhena/gogen/templates"
	"github.com/mirzaakhena/gogen/util"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type EntityModel struct {
	EntityName string
}

func NewEntityModel() (Commander, error) {
	flag.Parse()

	entityName := flag.Arg(1)
	if entityName == "" {
		return nil, fmt.Errorf("entity name must not empty. `gogen entity EntityName`")
	}

	return &EntityModel{
		EntityName: entityName,
	}, nil
}

func (obj *EntityModel) Run() error {

	err := util.CreateFolderIfNotExist("domain/entity")
	if err != nil {
		return err
	}

	{

		fset := token.NewFileSet()

		pkgs, err := parser.ParseDir(fset, "domain/entity", nil, parser.ParseComments)
		if err != nil {
			if err != nil {
				return err
			}
		}

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {

				for _, decl := range file.Decls {

					gen, ok := decl.(*ast.GenDecl)
					if !ok || gen.Tok != token.TYPE {
						continue
					}

					for _, specs := range gen.Specs {

						ts, ok := specs.(*ast.TypeSpec)
						if !ok {
							continue
						}

						if _, ok = ts.Type.(*ast.StructType); !ok {
							continue
						}

						// entity already exist, abort the command
						if ts.Name.String() == fmt.Sprintf("%s", obj.EntityName) {
							return nil
						}
					}
				}

			}
		}

	}

	{
		outputFile := fmt.Sprintf("domain/entity/%s.go", strings.ToLower(obj.EntityName))
		err = util.WriteFileIfNotExist(templates.EntityFile, outputFile, obj)
		if err != nil {
			return err
		}
	}

	return nil

}
