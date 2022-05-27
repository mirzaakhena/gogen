package gencontroller

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mirzaakhena/gogen/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
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

	if len(inputs) < 1 {

		err := fmt.Errorf("\n" +
			"   # Create a controller with gin as default web frameworkfor all usecases\n" +
			"   gogen controller restapi\n" +
			"     'restapi' is a gateway name\n" +
			"\n" +
			"   # Create a controller with with defined web framework for all usecases\n" +
			"   gogen controller restapi gin\n" +
			"     'restapi'     is a gateway name\n" +
			"     'CreateOrder' is an usecase name\n" +
			"\n" +
			"   # Create a controller with defined web framework and specific usecase\n" +
			"   gogen controller restapi gin CreateOrder\n" +
			"     'restapi'      is a gateway name\n" +
			"     'gin'          is a sample webframewrok. You may try the other one like: nethttp, echo, and gorilla\n" +
			"     'CreateOrder'  is an usecase name\n" +
			"\n")

		return err
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

		fileInfo, err := ioutil.ReadDir(usecaseFolderName)
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

	packagePath := utils.GetPackagePath()

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
		templateCode, err := getHandlerTemplate(obj.DriverName)
		if err != nil {
			return err
		}

		//templateWithData, err := utils.PrintTemplate(string(templateCode), obj)
		//if err != nil {
		//	return err
		//}

		filename := fmt.Sprintf("domain_%s/controller/%s/handler_%s.go", domainName, utils.LowerCase(controllerName), utils.LowerCase(usecase.Name))

		singleObj := ObjTemplateSingle{
			PackagePath:    obj.PackagePath,
			DomainName:     domainName,
			ControllerName: controllerName,
			DriverName:     driverName,
			Usecase:        usecase,
		}

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

	unexistedUsecases, err := getUnexistedUsecase(packagePath, domainName, controllerName, obj.Usecases)
	if err != nil {
		return err
	}

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
		{
			templateCode, err := getRouterInportTemplate(obj.DriverName)
			if err != nil {
				return err
			}

			templateWithData, err := utils.PrintTemplate(string(templateCode), singleObj)
			if err != nil {
				return err
			}

			dataInBytes, err := injectInportToStruct(obj, templateWithData)
			if err != nil {
				return err
			}

			// reformat router.go
			err = utils.Reformat(obj.getControllerRouterFileName(), dataInBytes)
			if err != nil {
				return err
			}
		}

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

	return nil

}

func injectUsecaseInportFields(usecaseFolderName string, usecaseName string, usecases []*Usecase) []*Usecase {

	inportRequestFields := make([]*StructField, 0)
	inportResponseFields := make([]*StructField, 0)
	fset := token.NewFileSet()
	utils.IsExist(fset, fmt.Sprintf("%s/%s", usecaseFolderName, usecaseName), func(file *ast.File, ts *ast.TypeSpec) bool {

		structObj, isStruct := ts.Type.(*ast.StructType)

		if isStruct {

			if utils.LowerCase(ts.Name.String()) == "inportrequest" {

				for _, f := range structObj.Fields.List {
					fieldType := utils.TypeHandler{PrefixExpression: utils.LowerCase(usecaseName)}.Mulai(f.Type)
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
					fieldType := utils.TypeHandler{PrefixExpression: utils.LowerCase(usecaseName)}.Mulai(f.Type)
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

		return false
	})

	return usecases
}

func getUnexistedUsecase(packagePath, domainName, controllerName string, currentUsecases []*Usecase) ([]*Usecase, error) {

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
	for _, decl := range astFile.Decls {

		if gen, ok := decl.(*ast.FuncDecl); ok {

			if gen.Recv == nil {
				continue
			}

			startExp, ok := gen.Recv.List[0].Type.(*ast.StarExpr)
			if !ok {
				continue
			}

			if startExp.X.(*ast.Ident).String() != "Controller" {
				continue
			}

			if gen.Name.String() != "RegisterRouter" {
				continue
			}

			routerLine = fset.Position(gen.Body.Rbrace).Line
			return routerLine, nil
		}

	}
	return 0, fmt.Errorf("register router Not found")
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
