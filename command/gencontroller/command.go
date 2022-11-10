package gencontroller

import (
	"bufio"
	"bytes"
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
	Usecases       []*Usecase
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

		frameworks := ""

		dirs, err := utils.AppTemplates.ReadDir("templates/controllers")
		if err != nil {
			return err
		}

		for i, dir := range dirs {

			name := dir.Name()

			if dir.IsDir() {

				if len(dirs) == 1 {
					frameworks = fmt.Sprintf("%s", frameworks)
					continue
				}

				if i == 0 {
					frameworks = fmt.Sprintf("%s,", name)
					continue
				}

				if i == len(dirs)-1 {
					frameworks = fmt.Sprintf("%s %s", frameworks, name)
					continue
				}

				frameworks = fmt.Sprintf("%s %s,", frameworks, name)
			}
		}

		msg := fmt.Errorf("\n"+
			"   # Create a controller for all usecases using gin as default web framework\n"+
			"   gogen controller restapi\n"+
			"     'restapi' is a gateway name\n"+
			"\n"+
			"   # Create a controller with for all usecases with selected framework\n"+
			"   You may try the other one like : %s\n"+
			"   in this example we use 'echo'\n"+
			"   gogen controller restapi echo\n"+
			"     'restapi'     is a gateway name\n"+
			"     'CreateOrder' is an usecase name\n"+
			"\n"+
			"   # Create a controller with defined web framework and specific usecase\n"+
			"   gogen controller restapi gin CreateOrder\n"+
			"     'restapi'      is a gateway name\n"+
			"     'gin'          is a sample webframework.\n"+
			"     'CreateOrder'  is an usecase name\n"+
			"\n", frameworks)

		return msg
	}

	domainName := utils.GetDefaultDomain()
	controllerName := inputs[0]

	driverName := "gin"
	if len(inputs) >= 2 {
		driverName = utils.LowerCase(inputs[1])
	}

	usecaseFolderName := fmt.Sprintf("domain_%s/usecase", domainName)

	usecaseNames := make([]string, 0)

	usecases := make([]*Usecase, 0)

	if len(inputs) >= 3 {
		usecaseNames = append(usecaseNames, inputs[2])

		usecaseName := utils.LowerCase(inputs[2])

		usecases = injectUsecaseInportFields(usecaseFolderName, usecaseName, usecases)

	} else {

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

	}

	obj := ObjTemplate{
		PackagePath:    packagePath,
		DomainName:     domainName,
		ControllerName: controllerName,
		DriverName:     driverName,
		Usecases:       usecases,
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	fileRenamer := map[string]string{
		"controllername": utils.LowerCase(controllerName),
		"domainname":     utils.LowerCase(domainName),
	}

	err = utils.CreateEverythingExactly("templates/controllers/", obj.DriverName, fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	// handler_<usecase>.go
	for _, usecase := range obj.Usecases {

		singleObj := ObjTemplateSingle{
			PackagePath:    obj.PackagePath,
			DomainName:     domainName,
			ControllerName: controllerName,
			DriverName:     driverName,
			Usecase:        usecase,
		}

		{
			templateCode, err := getHandlerTemplate(obj.DriverName)
			if err != nil {
				return err
			}

			//templateWithData, err := utils.PrintTemplate(string(templateCode), obj)
			//if err != nil {
			//	return err
			//}

			filename := fmt.Sprintf("domain_%s/controller/%s/handler_%s.go", domainName, utils.LowerCase(controllerName), utils.LowerCase(usecase.Name))

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

		if strings.HasPrefix(strings.ToLower(usecase.Name), "get") {

			templateCode, err := getHTTPClientGETTemplate(obj.DriverName)
			if err != nil {
				return err
			}

			filename := fmt.Sprintf("domain_%s/controller/%s/http_%s.http", domainName, utils.LowerCase(controllerName), utils.LowerCase(usecase.Name))

			_, err = utils.WriteFileIfNotExist(string(templateCode), filename, singleObj)
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(strings.ToLower(usecase.Name), "run") {

			templateCode, err := getHTTPClientPOSTTemplate(obj.DriverName)
			if err != nil {
				return err
			}

			filename := fmt.Sprintf("domain_%s/controller/%s/http_%s.http", domainName, utils.LowerCase(controllerName), utils.LowerCase(usecase.Name))

			_, err = utils.WriteFileIfNotExist(string(templateCode), filename, singleObj)
			if err != nil {
				return err
			}

		}

	}

	unexistedUsecases, err := getUnexistedUsecaseFromRouterBind(packagePath, domainName, controllerName, obj.Usecases)
	if err != nil {
		return err
	}

	//unexistedUsecases, err := getUnexistedUsecaseFromImport(packagePath, domainName, controllerName, obj.Usecases)
	//if err != nil {
	//	return err
	//}

	if len(unexistedUsecases) == 0 {
		// reformat router.go
		err = utils.Reformat(obj.getControllerRouterFileName(), nil)
		if err != nil {
			return err
		}
	}

	for _, usecase := range unexistedUsecases {

		singleObj := ObjTemplateSingle{
			PackagePath:    obj.PackagePath,
			DomainName:     domainName,
			ControllerName: controllerName,
			Usecase:        usecase,
			DriverName:     driverName,
		}

		//inject inport to struct
		//type Controller struct {
		//  Router            gin.IRouter
		//  CreateOrderInport createorder.Inport <----- here
		//}
		//{
		//	templateCode, err := getRouterInportTemplate(obj.DriverName)
		//	if err != nil {
		//		return err
		//	}
		//
		//	templateWithData, err := utils.PrintTemplate(string(templateCode), singleObj)
		//	if err != nil {
		//		return err
		//	}
		//
		//	dataInBytes, err := injectInportToStruct(obj, templateWithData)
		//	if err != nil {
		//		return err
		//	}
		//
		//	// reformat router.go
		//	err = utils.Reformat(obj.getControllerRouterFileName(), dataInBytes)
		//	if err != nil {
		//		return err
		//	}
		//}

		// inject router for register
		//func (r *Controller) RegisterRouter() {
		//  r.Router.POST("/createorder", r.authorized(), r.createOrderHandler(r.CreateOrderInport)) <-- here
		//}
		{
			templateCode, err := getRouterRegisterTemplate(obj.DriverName, usecase.Name)

			templateWithData, err := utils.PrintTemplate(string(templateCode), singleObj)
			if err != nil {
				return err
			}

			dataInBytes, err := injectRouterBind(obj, templateWithData)
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

	//fmt.Printf("Code is generated successfuly. If you saw the error like :\n")
	//fmt.Printf("- could not import \"some/path\" (no required module provides package \"some/path\")\n")
	//fmt.Printf("- unresolved type 'x'\n")
	//fmt.Printf("- cannot resolve symbol 'x'\n")
	//fmt.Printf("try to run the 'go mod tidy' manually in order to download required dependency\n")

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

	//usecases = append(usecases, &Usecase{
	//	Name:                 usecaseNameFromConst,
	//	InportRequestFields:  inportRequestFields,
	//	InportResponseFields: inportResponseFields,
	//})

	//utils.IsExist(fset, fmt.Sprintf("%s/%s", usecaseFolderName, usecaseName), func(file *ast.File, ts *ast.TypeSpec) bool {
	//
	//	structObj, isStruct := ts.Type.(*ast.StructType)
	//
	//	if isStruct {
	//
	//		if utils.LowerCase(ts.Name.String()) == "inportrequest" {
	//
	//			for _, f := range structObj.Fields.List {
	//				fieldType := utils.TypeHandler{PrefixExpression: utils.LowerCase(usecaseName)}.Start(f.Type)
	//				for _, name := range f.Names {
	//					inportRequestFields = append(inportRequestFields, &StructField{
	//						Name: name.String(),
	//						Type: fieldType,
	//					})
	//				}
	//			}
	//		}
	//
	//		if utils.LowerCase(ts.Name.String()) == "inportresponse" {
	//
	//			for _, f := range structObj.Fields.List {
	//				fieldType := utils.TypeHandler{PrefixExpression: utils.LowerCase(usecaseName)}.Start(f.Type)
	//				for _, name := range f.Names {
	//					inportResponseFields = append(inportResponseFields, &StructField{
	//						Name: name.String(),
	//						Type: fieldType,
	//					})
	//				}
	//			}
	//		}
	//
	//		if utils.LowerCase(ts.Name.String()) == fmt.Sprintf("%sinteractor", file.Name) {
	//			usecaseNameWithInteractor := ts.Name.String()
	//			usecaseNameOnly := usecaseNameWithInteractor[:strings.LastIndex(usecaseNameWithInteractor, "Interactor")]
	//			usecases = append(usecases, &Usecase{
	//				Name:                 usecaseNameOnly,
	//				InportRequestFields:  inportRequestFields,
	//				InportResponseFields: inportResponseFields,
	//			})
	//
	//		}
	//	}
	//
	//	return false
	//})

	return usecases
}

func getUnexistedUsecaseFromImport(packagePath, domainName, controllerName string, currentUsecases []*Usecase) ([]*Usecase, error) {

	unexistUsecase := make([]*Usecase, 0)

	routerFile := fmt.Sprintf("domain_%s/controller/%s/router.go", domainName, controllerName)

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, routerFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	mapUsecase := map[string]int{}

	// loop the outport for imports
	for _, decl := range astFile.Decls {

		if gen, ok := decl.(*ast.GenDecl); ok {

			if gen.Tok != token.IMPORT {
				continue
			}

			for _, spec := range gen.Specs {

				importSpec := spec.(*ast.ImportSpec)

				if strings.HasPrefix(importSpec.Path.Value, fmt.Sprintf("\"%s/domain_%s/usecase/", packagePath, domainName)) {
					readUsecase := importSpec.Path.Value[strings.LastIndex(importSpec.Path.Value, "/")+1:]
					uc := readUsecase[:len(readUsecase)-1]

					mapUsecase[uc] = 1

				}
			}

		}
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

		//if gen, ok := decl.(*ast.GenDecl); ok {
		//
		//	if gen.Tok != token.IMPORT {
		//		continue
		//	}
		//
		//	for _, spec := range gen.Specs {
		//
		//		importSpec := spec.(*ast.ImportSpec)
		//
		//		if strings.HasPrefix(importSpec.Path.Value, fmt.Sprintf("\"%s/domain_%s/usecase/", packagePath, domainName)) {
		//			readUsecase := importSpec.Path.Value[strings.LastIndex(importSpec.Path.Value, "/")+1:]
		//			uc := readUsecase[:len(readUsecase)-1]
		//
		//			mapUsecase[uc] = 1
		//
		//		}
		//	}
		//
		//}
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

func injectInportToStruct(obj ObjTemplate, templateWithData string) ([]byte, error) {

	inportLine, err := getInportLine(obj)
	if err != nil {
		return nil, err
	}

	controllerFile := obj.getControllerRouterFileName()

	file, err := os.Open(controllerFile)
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
	line := 0
	for scanner.Scan() {
		row := scanner.Text()

		if line == inportLine-1 {
			buffer.WriteString(templateWithData)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
		line++
	}

	return buffer.Bytes(), nil
}

func getInportLine(obj ObjTemplate) (int, error) {

	controllerFile := obj.getControllerRouterFileName()

	inportLine := 0
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, controllerFile, nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}

	// loop the outport for imports
	for _, decl := range astFile.Decls {

		if gen, ok := decl.(*ast.GenDecl); ok {

			if gen.Tok != token.TYPE {
				continue
			}

			for _, specs := range gen.Specs {

				ts, ok := specs.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if iStruct, ok := ts.Type.(*ast.StructType); ok {

					// check the specific struct name
					if ts.Name.String() != "Controller" {
						continue
					}

					inportLine = fset.Position(iStruct.Fields.Closing).Line
					return inportLine, nil
				}

			}

		}

	}

	return 0, fmt.Errorf(" Controller struct not found")

}

func getBindRouterLine(obj ObjTemplate) (int, error) {

	controllerFile := obj.getControllerRouterFileName()

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, controllerFile, nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}
	routerLine := 0

	//ast.Print(fset, astFile)

	ast.Inspect(astFile, func(node ast.Node) bool {

		for {
			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok {
				break
			}

			if funcDecl.Name.String() != "RegisterRouter" {
				break
			}

			routerLine = fset.Position(funcDecl.Body.Rbrace).Line

			return false
		}

		return true
	})

	if routerLine == 0 {
		return 0, fmt.Errorf("register router Not found")
	}

	return routerLine, nil

	//for _, decl := range astFile.Decls {
	//
	//	if gen, ok := decl.(*ast.FuncDecl); ok {
	//
	//		if gen.Recv == nil {
	//			continue
	//		}
	//
	//		startExp, ok := gen.Recv.List[0].Type.(*ast.StarExpr)
	//		if !ok {
	//			continue
	//		}
	//
	//		if startExp.X.(*ast.Ident).String() != "Controller" {
	//			continue
	//		}
	//
	//		if gen.Name.String() != "RegisterRouter" {
	//			continue
	//		}
	//
	//		routerLine = fset.Position(gen.Body.Rbrace).Line
	//		return routerLine, nil
	//	}
	//
	//}

}

func injectRouterBind(obj ObjTemplate, templateWithData string) ([]byte, error) {

	controllerFile := obj.getControllerRouterFileName()

	routerLine, err := getBindRouterLine(obj)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(controllerFile)
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
	line := 0
	for scanner.Scan() {
		row := scanner.Text()

		if line == routerLine-1 {
			buffer.WriteString(templateWithData)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
		line++
	}

	return buffer.Bytes(), nil

}

// getControllerRouterFileName ...
func (o ObjTemplate) getControllerRouterFileName() string {
	return fmt.Sprintf("domain_%s/controller/%s/router.go", utils.LowerCase(o.DomainName), utils.LowerCase(o.ControllerName))
}

func getHTTPClientGETTemplate(driverName string) ([]byte, error) {
	path := fmt.Sprintf("templates/controllers/%s/domain_${domainname}/controller/${controllername}/~httpclient-get._http", driverName)
	return utils.AppTemplates.ReadFile(path)
}

func getHTTPClientPOSTTemplate(driverName string) ([]byte, error) {
	path := fmt.Sprintf("templates/controllers/%s/domain_${domainname}/controller/${controllername}/~httpclient-post._http", driverName)
	return utils.AppTemplates.ReadFile(path)
}

func getHandlerTemplate(driverName string) ([]byte, error) {
	path := fmt.Sprintf("templates/controllers/%s/domain_${domainname}/controller/${controllername}/~handler._go", driverName)
	return utils.AppTemplates.ReadFile(path)
}

func getRouterInportTemplate(driverName string) ([]byte, error) {
	path := fmt.Sprintf("templates/controllers/%s/domain_${domainname}/controller/${controllername}/~inject-router-inport._go", driverName)
	return utils.AppTemplates.ReadFile(path)
}

func getRouterRegisterTemplate(driverName, usecase string) ([]byte, error) {
	if strings.HasPrefix(strings.ToLower(usecase), "get") {
		path := fmt.Sprintf("templates/controllers/%s/domain_${domainname}/controller/${controllername}/~inject-router-register-get._go", driverName)
		return utils.AppTemplates.ReadFile(path)
	}

	path := fmt.Sprintf("templates/controllers/%s/domain_${domainname}/controller/${controllername}/~inject-router-register-post._go", driverName)
	return utils.AppTemplates.ReadFile(path)
}
