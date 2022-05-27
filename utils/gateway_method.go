package utils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type gatewayMethod struct {
	packagePath string
}

var existingFunc map[string]int

func NewOutportMethodImpl(structName, gatewayRootFolderName string) (map[string]int, error) {

	existingFunc = map[string]int{}

	gm := gatewayMethod{packagePath: GetPackagePath()}

	err := gm.readStruct(structName, gatewayRootFolderName)
	if err != nil {
		return nil, err
	}

	return existingFunc, nil
}

// existingFunc map[string]int dibuang dari parameter
func (obj *gatewayMethod) readStruct(structName, folderPath string) error {

	//fmt.Printf(">>>>> structname: %s, folderpath:%s\n", structName, folderPath)

	if folderPath == "" {
		return nil
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, folderPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {

		// read file by file
		for _, file := range pkg.Files {

			importPaths := map[string]string{}

			for _, decl := range file.Decls {

				//ast.Print(fset, decl)

				switch gd := decl.(type) {

				case *ast.GenDecl:
					err := obj.generalDecl(structName, gd, folderPath, importPaths)
					if err != nil {
						return err
					}

				case *ast.FuncDecl:
					if !obj.findAndCollectImplMethod(gd, structName) {
						continue
					}
				}

			}

		}

	}

	return nil
}

func (obj *gatewayMethod) generalDecl(structName string, gd *ast.GenDecl, folderPath string, importPaths map[string]string) error {
	for _, spec := range gd.Specs {

		// handle import
		is, ok := spec.(*ast.ImportSpec)
		if ok {
			handleImports(obj.packagePath, is, importPaths)

			//fmt.Printf(">>>> %v\n", importPaths)
		}

		// it is type declaration
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		// the struct name must have a 'gateway' suffix
		if !strings.HasSuffix(strings.ToLower(ts.Name.String()), structName) {
			continue
		}

		// gateway must be a struct type
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			continue
		}

		// if struct list empty then nothing to do
		if st.Fields.List == nil {
			break
		}

		for _, fieldList := range st.Fields.List {

			err := obj.handleListFieldType(fieldList.Type, folderPath, importPaths)
			if err != nil {
				return err
			}

		}

	}
	return nil
}

func (obj *gatewayMethod) handleListFieldType(expr ast.Expr, folderPath string, importPaths map[string]string) error {
	switch ty := expr.(type) {

	case *ast.StarExpr: // can be ident or selector
		err := obj.handleListFieldType(ty.X, folderPath, importPaths)
		if err != nil {
			return err
		}

	case *ast.Ident: // still in the same package
		err := obj.readStruct(ty.Name, folderPath)
		if err != nil {
			return err
		}

	case *ast.SelectorExpr: // import from other package

		expression := ty.X.(*ast.Ident).String()
		pathWithGomod := importPaths[expression]
		pathOnly := strings.TrimPrefix(pathWithGomod, obj.packagePath+"/")
		structName := ty.Sel.String()
		err := obj.readStruct(structName, pathOnly)
		if err != nil {
			return err
		}

	}

	return nil
}

func (obj *gatewayMethod) findAndCollectImplMethod(fd *ast.FuncDecl, structName string) bool {
	if fd.Recv == nil {
		return false
	}

	switch ty := fd.Recv.List[0].Type.(type) {

	case *ast.Ident: // func (r domain) FindPayment() {}
		// read all the function that have receiver with gateway name
		if ty.String() != structName {
			return false
		}

	case *ast.StarExpr: // func (r *domain) FindPayment() {}
		// read all the function that have receiver with gateway name
		if ty.X.(*ast.Ident).String() != structName {
			return false
		}

	}

	// collect all existing function that have been there in the file
	existingFunc[fd.Name.String()] = 1

	return true
}
