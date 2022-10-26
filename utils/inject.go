package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
	"strings"
)

const injectedCodeLocation = "//!"

// InjectToInteractor ...
func InjectToInteractor(interactorFilename, injectedCode string) ([]byte, error) {

	existingFile := interactorFilename

	// open interactor file
	file, err := os.Open(existingFile)
	if err != nil {
		return nil, err
	}

	needToInject := false

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		// check the injected code in interactor
		if strings.TrimSpace(row) == injectedCodeLocation {

			needToInject = true

			// inject code
			buffer.WriteString(injectedCode)
			buffer.WriteString("\n")

			continue
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	// if no injected marker found, then abort the next step
	if !needToInject {
		return nil, nil
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	// rewrite the file
	if err := os.WriteFile(existingFile, buffer.Bytes(), 0644); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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

func InjectToMain(fset *token.FileSet, applicationName string) {

	//astFile, err := parser.ParseFile(fset, "main.go", nil, parser.ParseComments)
	//if err != nil {
	//	fmt.Printf("%v\n", err.Error())
	//	os.Exit(1)
	//}

	// ast.Print(fset, astFile)

	//pkgs, err := parser.ParseDir(fset, rootFolderName, nil, parser.ParseComments)
	//if err != nil {
	//	fmt.Printf("%v\n", err.Error())
	//	os.Exit(1)
	//}

	// in every declaration like type, func, const
	//for _, decl := range astFile.Decls {
	//
	//	// focus only to type
	//	funcDecl, ok := decl.(*ast.FuncDecl)
	//	if !ok || funcDecl.Name.String() != "main" {
	//		continue
	//	}
	//
	//	for _, stmts := range funcDecl.Body.List {
	//
	//		assignStmt, ok := stmts.(*ast.AssignStmt)
	//		if !ok {
	//			continue
	//		}
	//
	//		foundRegistryContract := false
	//		for _, rhs := range assignStmt.Rhs {
	//
	//			compositeLit, ok := rhs.(*ast.CompositeLit)
	//			if !ok {
	//				continue
	//			}
	//
	//			// map dengan key string
	//			mapType := compositeLit.Type.(*ast.MapType)
	//			if mapType.Key.(*ast.Ident).String() != "string" {
	//				continue
	//			}
	//
	//			// map dengan value func
	//			funcType, ok := mapType.Value.(*ast.FuncType)
	//			if !ok {
	//				continue
	//			}
	//
	//			// func yang mereturn RegistryContract
	//			for _, resultField := range funcType.Results.List {
	//				selectorExpr := resultField.Type.(*ast.SelectorExpr)
	//
	//				if selectorExpr.Sel.String() == "RegistryContract" {
	//					foundRegistryContract = true
	//					break
	//				}
	//
	//			}
	//
	//			if !foundRegistryContract {
	//				continue
	//			}
	//
	//			//ast.Print(fset, compositeLit)
	//
	//			var valuePos token.Pos
	//			if len(compositeLit.Elts) == 0 {
	//				valuePos = compositeLit.Lbrace + 3
	//			} else {
	//				valuePos = compositeLit.Elts[len(compositeLit.Elts)-1].End() + 4
	//			}
	//
	//			for _, elt := range compositeLit.Elts {
	//				kvExpr := elt.(*ast.KeyValueExpr)
	//				callExpr := kvExpr.Value.(*ast.CallExpr)
	//				selectorExpr := callExpr.Fun.(*ast.SelectorExpr)
	//				if selectorExpr.Sel.String() == fmt.Sprintf("New%s", PascalCase(applicationName)) {
	//					return
	//				}
	//			}
	//
	//			compositeLit.Elts = append(compositeLit.Elts, &ast.KeyValueExpr{
	//				Key: &ast.BasicLit{
	//					Kind:     token.STRING,
	//					ValuePos: valuePos,
	//					Value:    fmt.Sprintf("\"%s\"", LowerCase(applicationName)),
	//				},
	//				Value: &ast.CallExpr{
	//					Fun: &ast.SelectorExpr{
	//						X:   &ast.Ident{Name: "application"},
	//						Sel: &ast.Ident{Name: fmt.Sprintf("New%s", PascalCase(applicationName))},
	//					},
	//				},
	//				Colon: valuePos,
	//			})
	//
	//			compositeLit.Rbrace = valuePos + 2
	//
	//			//ast.Print(fset, compositeLit)
	//
	//			{
	//				f, err := os.Create("main.go")
	//				if err != nil {
	//					return
	//				}
	//				//defer func(f *os.File) {
	//				//	err := f.Close()
	//				//	if err != nil {
	//				//		os.Exit(1)
	//				//	}
	//				//}(f)
	//				err = printer.Fprint(f, fset, astFile)
	//				if err != nil {
	//					return
	//				}
	//
	//				err = f.Close()
	//				if err != nil {
	//					os.Exit(1)
	//				}
	//			}
	//			//{
	//			//	err := printer.Fprint(os.Stdout, fset, astFile)
	//			//	if err != nil {
	//			//		return
	//			//	}
	//			//}
	//		}
	//	}
	//}
}

func InjectToErrorEnum(fset *token.FileSet, filepath string, errorName, separator string) {

	astFile, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	// in every declaration like type, func, const
	for _, decl := range astFile.Decls {

		genDecl := decl.(*ast.GenDecl)

		if genDecl.Tok != token.CONST {
			continue
		}

		var errorNumber = 1000
		for _, spec := range genDecl.Specs {

			valueSpec := spec.(*ast.ValueSpec)
			if len(valueSpec.Values) == 0 {
				break
			}

			if len(valueSpec.Names) > 0 && strings.ToLower(valueSpec.Names[0].String()) == strings.ToLower(errorName) {
				fmt.Printf("error code %s already exist\n", errorName)
				return
			}

			basicList := valueSpec.Values[0].(*ast.BasicLit)
			errorCodeWithMessage := strings.Split(basicList.Value, " ")
			if len(errorCodeWithMessage) == 0 {
				continue
			}

			errorCodeOnly := strings.Split(errorCodeWithMessage[0], separator)
			if len(errorCodeOnly) < 2 || errorCodeOnly[1] == "" {
				continue
			}

			n, err := strconv.Atoi(errorCodeOnly[1])
			if err != nil {
				continue
			}
			errorNumber = n
		}

		errorValue := fmt.Sprintf("\"%s%04d %s\"", separator, errorNumber+1, SpaceCase(errorName))

		genDecl.Specs = append(genDecl.Specs, &ast.ValueSpec{
			Names:  []*ast.Ident{{Name: PascalCase(errorName)}},
			Type:   &ast.SelectorExpr{X: &ast.Ident{Name: "apperror"}, Sel: &ast.Ident{Name: "ErrorType"}},
			Values: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: errorValue}},
		})

		//ast.Print(fset, decl)

	}

	{
		f, err := os.Create(filepath)
		if err != nil {
			return
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				os.Exit(1)
			}
		}(f)
		err = printer.Fprint(f, fset, astFile)
		if err != nil {
			return
		}
	}

	//{
	//	err := printer.Fprint(os.Stdout, fset, astFile)
	//	if err != nil {
	//		return
	//	}
	//}
}
