package entity

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const injectedCodeLocation = "//!"

// GetPackageName ...
func GetPackageName(rootFolderName string) string {
	i := strings.LastIndex(rootFolderName, "/")
	return rootFolderName[i+1:]
}

func InjectCodeAtTheEndOfFile(filename, templateCode string) ([]byte, error) {

	// reopen the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	// write the template in the end of file
	buffer.WriteString(templateCode)
	buffer.WriteString("\n")

	return buffer.Bytes(), nil

}

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

					// focus only to some struct for example
					//if _, ok = ts.Type.(*ast.StructType); !ok {
					//  continue
					//}

					if foundExactType(file, ts) {
						return true
					}

					// entity already exist, abort the command
					//if ts.Name.String() != typeName {
					//  continue
					//}
					//
					//fmt.Printf("File : %v\n", fset.Position(file.Pos()).Filename)

				}
			}
		}
	}
	return false
}
