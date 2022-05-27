package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func IsExist(fset *token.FileSet, rootFolderName string, foundExactType func(file *ast.File, ts *ast.TypeSpec) bool) bool {

	pkgs, err := parser.ParseDir(fset, rootFolderName, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	// in every package
	for _, pkg := range pkgs {

		// in every files
		for _, file := range pkg.Files {

			// in every declaration like type, func, const
			for _, decl := range file.Decls {

				// focus only to type
				gen, ok := decl.(*ast.GenDecl)
				if !ok || gen.Tok != token.TYPE {
					continue
				}

				for _, specs := range gen.Specs {

					ts, ok := specs.(*ast.TypeSpec)
					if !ok {
						continue
					}

					if foundExactType(file, ts) {
						return true
					}

				}
			}
		}
	}
	return false
}
