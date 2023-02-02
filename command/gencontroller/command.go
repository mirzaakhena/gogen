package gencontroller

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath    string
	DomainName     string
	ControllerName string
	DriverName     string
	//Usecases       []*Usecase
}

// ObjTemplateSingle ...
type ObjTemplateSingle struct {
	PackagePath    string
	DomainName     string
	ControllerName string
	DriverName     string
	Usecase        *Usecase
}

type Usecase struct {
	Name                 string
	InportRequestFields  []*StructField
	InportResponseFields []*StructField
}

type StructField struct {
	Name string
	Type string
}

func Run(inputs ...string) error {

	packagePath := utils.GetPackagePath()

	if len(inputs) < 1 {

		msg := fmt.Errorf("\n" +
			"   # Create a controller for all usecases using default as default web framework\n" +
			"   gogen controller restapi\n" +
			"     'restapi' is a gateway name" +
			"\n")

		return msg
	}

	gcfg := utils.GetGogenConfig()
	controllerName := inputs[0]
	driverName := gcfg.Controller

	usecaseFolderName := fmt.Sprintf("domain_%s/usecase", gcfg.Domain)

	usecases := make([]*Usecase, 0)

	err := utils.CreateEverythingExactly("templates/", "shared", nil, nil, utils.AppTemplates)
	if err != nil {
		return err
	}

	fileInfo, err := os.ReadDir(usecaseFolderName)
	if err != nil {
		return err
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			continue
		}
		usecases = injectUsecaseInportFields(usecaseFolderName, file.Name(), usecases)

	}

	//}

	// siapkan obj controller yg akan di inject-kan, berisi list of usecases
	obj := ObjTemplate{
		PackagePath:    packagePath,
		DomainName:     gcfg.Domain,
		ControllerName: controllerName,
		DriverName:     driverName,
		//Usecases:       usecases,
	}

	fileRenamer := map[string]string{
		"controllername": utils.LowerCase(controllerName),
		"domainname":     utils.LowerCase(gcfg.Domain),
	}

	err = utils.CreateEverythingExactly2(".gogen/templates/controller/", obj.DriverName, fileRenamer, obj)
	if err != nil {
		return err
	}

	// -------------- done --------------

	// handler_<usecase>.go
	for _, usecase := range usecases {

		singleObj := ObjTemplateSingle{
			PackagePath:    obj.PackagePath,
			DomainName:     gcfg.Domain,
			ControllerName: controllerName,
			DriverName:     driverName,
			Usecase:        usecase,
		}

		// khusus nge-handle handler_usecasename.go saja
		{
			templateCode, err := getHandlerTemplate(obj.DriverName)
			if err != nil {
				return err
			}

			//templateWithData, err := utils.PrintTemplate(string(templateCode), obj)
			//if err != nil {
			//	return err
			//}

			filename := fmt.Sprintf("domain_%s/controller/%s/handler_%s.go", gcfg.Domain, utils.LowerCase(controllerName), utils.LowerCase(usecase.Name))

			_, err = utils.WriteFileIfNotExist(string(templateCode), filename, singleObj)
			if err != nil {
				return err
			}

			// reformat router.go
			err = utils.Reformat(filename, nil)
			if err != nil {
				return err
			}
		}

		// khusus handle httpclient.http saja
		{

			if strings.HasPrefix(strings.ToLower(usecase.Name), "get") {

				templateCode, err := getHTTPClientTemplate(obj.DriverName, "get")
				if err != nil {
					return err
				}

				filename := fmt.Sprintf("domain_%s/controller/%s/http_%s.http", gcfg.Domain, utils.LowerCase(controllerName), utils.LowerCase(usecase.Name))

				_, err = utils.WriteFileIfNotExist(string(templateCode), filename, singleObj)
				if err != nil {
					return err
				}

			} else if strings.HasPrefix(strings.ToLower(usecase.Name), "run") {

				templateCode, err := getHTTPClientTemplate(obj.DriverName, "post")
				if err != nil {
					return err
				}

				filename := fmt.Sprintf("domain_%s/controller/%s/http_%s.http", gcfg.Domain, utils.LowerCase(controllerName), utils.LowerCase(usecase.Name))

				_, err = utils.WriteFileIfNotExist(string(templateCode), filename, singleObj)
				if err != nil {
					return err
				}

			}

		}

	}

	// ambil semua usecase yang BELUM terdaftar di router
	unexistedUsecases, err := getUnexistedUsecaseFromRouterBind(packagePath, gcfg.Domain, controllerName, usecases)
	if err != nil {
		return err
	}

	//unexistedUsecases, err := getUnexistedUsecaseFromImport(packagePath, gcfg, controllerName, obj.Usecases)
	//if err != nil {
	//	return err
	//}

	//if len(unexistedUsecases) == 0 {
	//	// reformat router.go
	//	err = utils.Reformat(obj.getControllerRouterFileName(), nil)
	//	if err != nil {
	//		return err
	//	}
	//}

	for _, usecase := range unexistedUsecases {

		singleObj := ObjTemplateSingle{
			PackagePath:    obj.PackagePath,
			DomainName:     gcfg.Domain,
			ControllerName: controllerName,
			Usecase:        usecase,
			DriverName:     driverName,
		}

		// inject router for register
		//func (r *Controller) RegisterRouter() {
		//  r.Router.POST("/createorder", r.authorized(), r.createOrderHandler(r.CreateOrderInport)) <-- here
		//}
		{

			templateCode, err := getRouterRegisterTemplate(obj.DriverName, usecase.Name)

			templateHasBeenInjected, err := utils.PrintTemplate(string(templateCode), singleObj)
			if err != nil {
				return err
			}

			//dataInBytes, err := injectRouterBind(obj, templateHasBeenInjected)
			//if err != nil {
			//	return err
			//}

			controllerFile := obj.getControllerRouterFileName()
			dataInBytes, err := utils.InjectToCode(controllerFile, templateHasBeenInjected)
			if err != nil {
				return err
			}

			// reformat router.go
			err = utils.Reformat(obj.getControllerRouterFileName(), dataInBytes)
			if err != nil {
				return err
			}
		}

	}

	return nil

}

func injectUsecaseInportFields(usecaseFolderName string, usecaseName string, usecases []*Usecase) []*Usecase {

	inportRequestFields := make([]*StructField, 0)
	inportResponseFields := make([]*StructField, 0)
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, fmt.Sprintf("%s/%s", usecaseFolderName, usecaseName), nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	//usecaseNameFromConst := ""

	// in every package
	for _, pkg := range pkgs {

		// in every files
		for _, file := range pkg.Files {

			// in every declaration like type, func, const
			for _, decl := range file.Decls {

				// focus only to type
				gen, ok := decl.(*ast.GenDecl)

				//if ok && gen.Tok == token.CONST {
				//
				//	for _, spec := range gen.Specs {
				//
				//		valueSpec := spec.(*ast.ValueSpec)
				//		for _, name := range valueSpec.Names {
				//			if name.Name != "Name" {
				//				break
				//			}
				//		}
				//		for _, v := range valueSpec.Values {
				//			bl := v.(*ast.BasicLit)
				//			usecaseNameFromConst = strings.ReplaceAll(bl.Value, "\"", "")
				//		}
				//
				//	}
				//
				//} else

				if ok && gen.Tok == token.TYPE {

					for _, specs := range gen.Specs {

						ts, ok := specs.(*ast.TypeSpec)
						if !ok {
							continue
						}

						structObj, isStruct := ts.Type.(*ast.StructType)

						if isStruct {

							if utils.LowerCase(ts.Name.String()) == "inportrequest" {

								for _, f := range structObj.Fields.List {
									fieldType := utils.TypeHandler{PrefixExpression: utils.LowerCase(usecaseName)}.Start(f.Type)
									for _, name := range f.Names {
										inportRequestFields = append(inportRequestFields, &StructField{
											Name: name.String(),
											Type: fieldType,
										})
									}
								}
							}

							if utils.LowerCase(ts.Name.String()) == "inportresponse" {

								for _, f := range structObj.Fields.List {
									fieldType := utils.TypeHandler{PrefixExpression: utils.LowerCase(usecaseName)}.Start(f.Type)
									for _, name := range f.Names {
										inportResponseFields = append(inportResponseFields, &StructField{
											Name: name.String(),
											Type: fieldType,
										})
									}
								}
							}

							if utils.LowerCase(ts.Name.String()) == fmt.Sprintf("%sinteractor", file.Name) {
								usecaseNameWithInteractor := ts.Name.String()
								usecaseNameOnly := usecaseNameWithInteractor[:strings.LastIndex(usecaseNameWithInteractor, "Interactor")]
								usecases = append(usecases, &Usecase{
									Name:                 usecaseNameOnly,
									InportRequestFields:  inportRequestFields,
									InportResponseFields: inportResponseFields,
								})

							}
						}

					}

				}
			}
		}
	}

	return usecases
}

const HandlerSuffix = "handler"

func getUnexistedUsecaseFromRouterBind(packagePath, domainName, controllerName string, currentUsecases []*Usecase) ([]*Usecase, error) {

	unexistUsecase := make([]*Usecase, 0)

	routerFile := fmt.Sprintf("domain_%s/controller/%s/router.go", domainName, controllerName)

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, routerFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	//ast.Print(fset, astFile)

	mapUsecase := map[string]int{}

	// loop the outport for imports
	for _, decl := range astFile.Decls {

		ast.Inspect(decl, func(n ast.Node) bool {
			var s string
			switch x := n.(type) {
			case *ast.CallExpr:
				z, ok := x.Fun.(*ast.SelectorExpr)
				if ok {
					s = z.Sel.Name
				}
			}

			lowerCaseUsecase := strings.ToLower(s)
			if strings.HasSuffix(lowerCaseUsecase, HandlerSuffix) {
				uc := lowerCaseUsecase[:strings.LastIndex(lowerCaseUsecase, HandlerSuffix)]
				mapUsecase[uc] = 1
			}
			return true
		})
	}

	for _, usecase := range currentUsecases {
		_, exist := mapUsecase[utils.LowerCase(usecase.Name)]
		if exist {
			continue
		}

		unexistUsecase = append(unexistUsecase, usecase)
	}

	return unexistUsecase, nil
}

//func getBindRouterLine(obj ObjTemplate) (int, error) {
//
//	controllerFile := obj.getControllerRouterFileName()
//
//	fset := token.NewFileSet()
//	astFile, err := parser.ParseFile(fset, controllerFile, nil, parser.ParseComments)
//	if err != nil {
//		return 0, err
//	}
//	routerLine := 0
//
//	//ast.Print(fset, astFile)
//
//	ast.Inspect(astFile, func(node ast.Node) bool {
//
//		for {
//			funcDecl, ok := node.(*ast.FuncDecl)
//			if !ok {
//				break
//			}
//
//			if funcDecl.Name.String() != "RegisterRouter" {
//				break
//			}
//
//			routerLine = fset.Position(funcDecl.Body.Rbrace).Line
//
//			return false
//		}
//
//		return true
//	})
//
//	if routerLine == 0 {
//		return 0, fmt.Errorf("register router Not found")
//	}
//
//	return routerLine, nil
//
//}

//func injectRouterBind(obj ObjTemplate, templateWithData string) ([]byte, error) {
//
//	controllerFile := obj.getControllerRouterFileName()
//
//	routerLine, err := getBindRouterLine(obj)
//	if err != nil {
//		return nil, err
//	}
//
//	file, err := os.Open(controllerFile)
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if err := file.Close(); err != nil {
//			return
//		}
//	}()
//
//	scanner := bufio.NewScanner(file)
//	var buffer bytes.Buffer
//	line := 0
//	for scanner.Scan() {
//		row := scanner.Text()
//
//		if line == routerLine-1 {
//			buffer.WriteString(templateWithData)
//			buffer.WriteString("\n")
//		}
//
//		buffer.WriteString(row)
//		buffer.WriteString("\n")
//		line++
//	}
//
//	return buffer.Bytes(), nil
//
//}

// getControllerRouterFileName ...
func (o ObjTemplate) getControllerRouterFileName() string {
	return fmt.Sprintf("domain_%s/controller/%s/router.go", utils.LowerCase(o.DomainName), utils.LowerCase(o.ControllerName))
}

func getHTTPClientTemplate(driverName, method string) ([]byte, error) {
	path := fmt.Sprintf(".gogen/templates/controller/%s/domain_${domainname}/controller/${controllername}/~httpclient-%s._http", driverName, method)
	return os.ReadFile(path)
}

func getHandlerTemplate(driverName string) ([]byte, error) {
	path := fmt.Sprintf(".gogen/templates/controller/%s/domain_${domainname}/controller/${controllername}/~handler._go", driverName)
	return os.ReadFile(path)
}

func getRouterRegisterTemplate(driverName, usecase string) ([]byte, error) {
	method := "post"
	if strings.HasPrefix(strings.ToLower(usecase), "get") {
		method = "get"
	}
	path := fmt.Sprintf(".gogen/templates/controller/%s/domain_${domainname}/controller/${controllername}/~inject-router-register-%s._go", driverName, method)
	return os.ReadFile(path)
}

//func injectInportToStruct(obj ObjTemplate, templateWithData string) ([]byte, error) {
//
//	inportLine, err := getInportLine(obj)
//	if err != nil {
//		return nil, err
//	}
//
//	controllerFile := obj.getControllerRouterFileName()
//
//	file, err := os.Open(controllerFile)
//	if err != nil {
//		return nil, err
//	}
//
//	defer func() {
//		if err := file.Close(); err != nil {
//			return
//		}
//	}()
//
//	scanner := bufio.NewScanner(file)
//	var buffer bytes.Buffer
//	line := 0
//	for scanner.Scan() {
//		row := scanner.Text()
//
//		if line == inportLine-1 {
//			buffer.WriteString(templateWithData)
//			buffer.WriteString("\n")
//		}
//
//		buffer.WriteString(row)
//		buffer.WriteString("\n")
//		line++
//	}
//
//	return buffer.Bytes(), nil
//}

//func getInportLine(obj ObjTemplate) (int, error) {
//
//	controllerFile := obj.getControllerRouterFileName()
//
//	inportLine := 0
//	fset := token.NewFileSet()
//	astFile, err := parser.ParseFile(fset, controllerFile, nil, parser.ParseComments)
//	if err != nil {
//		return 0, err
//	}
//
//	// loop the outport for imports
//	for _, decl := range astFile.Decls {
//
//		if gen, ok := decl.(*ast.GenDecl); ok {
//
//			if gen.Tok != token.TYPE {
//				continue
//			}
//
//			for _, specs := range gen.Specs {
//
//				ts, ok := specs.(*ast.TypeSpec)
//				if !ok {
//					continue
//				}
//
//				if iStruct, ok := ts.Type.(*ast.StructType); ok {
//
//					// check the specific struct name
//					if ts.Name.String() != "Controller" {
//						continue
//					}
//
//					inportLine = fset.Position(iStruct.Fields.Closing).Line
//					return inportLine, nil
//				}
//
//			}
//
//		}
//
//	}
//
//	return 0, fmt.Errorf(" Controller struct not found")
//
//}
//func getUnexistedUsecaseFromImport(packagePath, domainName, controllerName string, currentUsecases []*Usecase) ([]*Usecase, error) {
//
//	unexistUsecase := make([]*Usecase, 0)
//
//	routerFile := fmt.Sprintf("domain_%s/controller/%s/router.go", domainName, controllerName)
//
//	fset := token.NewFileSet()
//	astFile, err := parser.ParseFile(fset, routerFile, nil, parser.ParseComments)
//	if err != nil {
//		return nil, err
//	}
//
//	mapUsecase := map[string]int{}
//
//	// loop the outport for imports
//	for _, decl := range astFile.Decls {
//
//		if gen, ok := decl.(*ast.GenDecl); ok {
//
//			if gen.Tok != token.IMPORT {
//				continue
//			}
//
//			for _, spec := range gen.Specs {
//
//				importSpec := spec.(*ast.ImportSpec)
//
//				if strings.HasPrefix(importSpec.Path.Value, fmt.Sprintf("\"%s/domain_%s/usecase/", packagePath, domainName)) {
//					readUsecase := importSpec.Path.Value[strings.LastIndex(importSpec.Path.Value, "/")+1:]
//					uc := readUsecase[:len(readUsecase)-1]
//
//					mapUsecase[uc] = 1
//
//				}
//			}
//
//		}
//	}
//
//	for _, usecase := range currentUsecases {
//		_, exist := mapUsecase[utils.LowerCase(usecase.Name)]
//		if exist {
//			continue
//		}
//
//		unexistUsecase = append(unexistUsecase, usecase)
//	}
//
//	return unexistUsecase, nil
//}
