package genapplication

import (
	"bufio"
	"bytes"
	"encoding/json"
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
	PackagePath     string
	DomainName      string
	ApplicationName string
	ControllerName  *string
	GatewayName     *string
	UsecaseNames    []string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a application for all controller\n" +
			"   gogen application appone\n" +
			"     'appone'       is an application name\n" +
			"\n" +
			"   # Create a application for specific controller\n" +
			"   gogen application appone restapi\n" +
			"     'appone'       is an application name\n" +
			"     'restapi'      is a controller name\n" +
			"\n" +
			"   # Create a application for specific controller and gateway\n" +
			"   gogen application appone restapi prod\n" +
			"     'appone'       is an application name\n" +
			"     'restapi'      is a controller name\n" +
			"     'prod'         is a gateway name\n" +
			"\n")

		return err
	}

	packagePath := utils.GetPackagePath()
	domainName := utils.GetDefaultDomain()
	applicationName := inputs[0]

	obj := &ObjTemplate{
		PackagePath:     packagePath,
		DomainName:      domainName,
		ApplicationName: applicationName,
		ControllerName:  nil,
		GatewayName:     nil,
	}

	if len(inputs) >= 2 {
		obj.ControllerName = &inputs[1]
	}

	if len(inputs) >= 3 {
		obj.GatewayName = &inputs[2]
	}

	var objController, objGateway string

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	{ //---------

		driverName := "gin"

		// if controller name is not given, then we will do auto controller discovery strategy
		if obj.ControllerName == nil {

			// look up the controller by foldername
			objControllers, err := findAllObjController(domainName)
			if err != nil {
				return err
			}

			// if there is more than one controller
			if len(objControllers) > 1 {
				names := make([]string, 0)

				// collect all the controller name
				for _, g := range objControllers {
					names = append(names, g)
				}

				// return error
				return fmt.Errorf("select one of this controller %v", names)
			}

			// currently, we are expecting only one gateway
			objController = objControllers[0]

		} else {

			var err error

			// when controller name is given
			objController, err = findControllerByName(domainName, *obj.ControllerName)
			if err != nil {
				return err
			}

			// in case the controller name is not found
			if objController == "" {
				return fmt.Errorf("no controller with name %s found", *obj.ControllerName)
			}

		}

		// if gateway name is not given, then we will do auto gateway discovery strategy
		if obj.GatewayName == nil {

			// look up the gateway by foldername
			objGateways, err := findAllObjGateway(domainName)
			if err != nil {
				return err
			}

			// if there is more than one gateway
			if len(objGateways) > 1 {
				names := make([]string, 0)

				// collect all the gateway name
				for _, g := range objGateways {
					names = append(names, g)
				}

				// return error
				return fmt.Errorf("select one of this gateways %v", names)
			}

			// currently, we are expecting only one gateway
			objGateway = objGateways[0]

		} else {

			var err error

			// when gateway name is given
			objGateway, err = findGatewayByName(domainName, *obj.GatewayName)
			if err != nil {
				return err
			}

			// in case the gateway name is not found
			if objGateway == "" {
				return fmt.Errorf("no gateway with name %s found", *obj.GatewayName)
			}
		}

		// ================

		// get all usecase from controller
		usecaseNames, err := findAllUsecaseInportNameFromController(domainName, objController)
		if err != nil {
			return err
		}

		existingUsecaseNames := map[string]int{}

		appFile := fmt.Sprintf("application/app_%s.go", strings.ToLower(applicationName))

		appFileIsExist := false

		if utils.IsFileExist(appFile) {
			appFileIsExist = true
			// get all existing usecase from application
			existingUsecaseNames, err = findExistingUsecaseInApplication(appFile)
			if err != nil {
				return err
			}
		}

		unexistUsecase := make([]string, 0)

		for _, usecaseName := range usecaseNames {
			_, exist := existingUsecaseNames[utils.LowerCase(usecaseName)]
			if exist {
				continue
			}

			unexistUsecase = append(unexistUsecase, usecaseName)
		}

		obj.GatewayName = &objGateway
		obj.ControllerName = &objController
		obj.UsecaseNames = unexistUsecase

		fileRenamer := map[string]string{
			"domainname":      utils.LowerCase(domainName),
			"applicationname": utils.LowerCase(applicationName),
		}

		err = utils.CreateEverythingExactly("templates/application/", driverName, fileRenamer, obj, utils.AppTemplates)
		if err != nil {
			return err
		}

		if appFileIsExist {

			for _, usecase := range unexistUsecase {

				dataInBytes, err := InjectRegisterUsecaseInApplication(usecase, appFile)
				if err != nil {
					return err
				}

				// reformat router.go
				err = utils.Reformat(appFile, dataInBytes)
				if err != nil {
					return err
				}

			}

		}

	} //---------

	// inject to main.go
	{
		// TODO find existing main.go check the existing application

		dataInBytes, err := InjectApplicationInMain(applicationName)
		if err != nil {
			return err
		}

		if dataInBytes != nil {
			// reformat router.go
			err = utils.Reformat("main.go", dataInBytes)
			if err != nil {
				return err
			}
		}

	}

	//{
	//	fset := token.NewFileSet()
	//	utils.InjectToMain(fset, applicationName)
	//}

	// inject into config.json
	{
		contentInBytes, err := os.ReadFile("config.json")
		if err != nil {
			panic(err.Error())
		}

		var cfg map[string]any

		err = json.Unmarshal(contentInBytes, &cfg)
		if err != nil {
			panic(err.Error())
		}

		portStart := 8000

		cfgServers, exist := cfg["servers"]
		if !exist {
			cfg["servers"] = map[string]any{
				applicationName: map[string]any{
					"address": fmt.Sprintf(":%d", portStart),
				},
			}
		} else {
			n := len(cfgServers.(map[string]any))
			cfgServers.(map[string]any)[applicationName] = map[string]any{
				"address": fmt.Sprintf(":%d", portStart+n),
			}
		}

		arrBytes, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			panic(err.Error())
		}

		err = os.WriteFile("config.json", arrBytes, os.ModeAppend)
		if err != nil {
			panic(err)
		}

	}

	//{
	//	bytes, err := os.ReadFile("docker-compose.yml")
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//
	//	var dc map[string]any
	//
	//	err = yaml.Unmarshal(bytes, &dc)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//
	//	services, exist := dc["services"]
	//	if !exist {
	//		panic("services not exist")
	//	}
	//
	//	mapServices, ok := services.(map[string]any)
	//	if ok {
	//
	//		for x := range mapServices {
	//
	//			if x == applicationName {
	//				// we already have it then nothing to add
	//				return nil
	//			}
	//
	//		}
	//
	//	}
	//
	//	if mapServices == nil {
	//		mapServices = map[string]any{}
	//	}
	//
	//	mapServices[applicationName] = map[string]any{
	//		"container_name": applicationName,
	//		"build":          `.`,
	//		"entrypoint":     fmt.Sprintf("./%s %s", utils.GetExecutableName(), applicationName),
	//		"ports":          []string{"8080:8080"},
	//	}
	//
	//	arrBytes, err := yaml.Marshal(dc)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//
	//	err = os.WriteFile("docker-compose.yml", arrBytes, os.ModeAppend)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//}

	return nil

}

func findAllObjController(domainName string) ([]string, error) {

	controllerFolder := getControllerFolder(domainName)

	if !utils.IsFileExist(controllerFolder) {
		return nil, fmt.Errorf("controller is not created yet")
	}

	dir, err := os.ReadDir(controllerFolder)
	if err != nil {
		return nil, err
	}

	controllers := make([]string, 0)

	for _, d := range dir {
		if !d.IsDir() {
			continue
		}

		g, err := findControllerByName(domainName, d.Name())
		if err != nil {
			return nil, err
		}
		if g == "" {
			continue
		}

		controllers = append(controllers, g)

	}

	return controllers, nil
}

func findControllerByName(domainName, controllerName string) (string, error) {
	folderName := fmt.Sprintf("%s/%s", getControllerFolder(domainName), strings.ToLower(controllerName))

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, folderName, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	for _, pkg := range pkgs {

		// read file by file
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

					if _, ok := ts.Type.(*ast.StructType); ok {

						// check the specific struct name
						if !strings.HasSuffix(strings.ToLower(ts.Name.String()), "controller") {
							continue
						}

						return pkg.Name, nil

						//inportLine = fset.Position(iStruct.Fields.Closing).Line
						//return inportLine, nil
					}
				}

			}

		}

	}

	return "", nil
}

func findAllObjGateway(domainName string) ([]string, error) {

	gatewayFolder := getGatewayFolder(domainName)

	if !utils.IsFileExist(gatewayFolder) {
		return nil, fmt.Errorf("gateway is not created yet")
	}

	dir, err := os.ReadDir(gatewayFolder)
	if err != nil {
		return nil, err
	}

	gateways := make([]string, 0)

	for _, d := range dir {
		if !d.IsDir() {
			continue
		}

		g, err := findGatewayByName(domainName, d.Name())
		if err != nil {
			return nil, err
		}
		if g == "" {
			continue
		}

		gateways = append(gateways, g)

	}

	return gateways, nil
}

func findGatewayByName(domainName, gatewayName string) (string, error) {
	folderName := fmt.Sprintf("%s/%s", getGatewayFolder(domainName), strings.ToLower(gatewayName))

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, folderName, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	for _, pkg := range pkgs {

		// read file by file
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

					if _, ok := ts.Type.(*ast.StructType); ok {

						// check the specific struct name
						if !strings.HasSuffix(strings.ToLower(ts.Name.String()), "gateway") {
							continue
						}

						return pkg.Name, nil

						//inportLine = fset.Position(iStruct.Fields.Closing).Line
						//return inportLine, nil
					}
				}

			}

		}

	}

	return "", nil
}

func findExistingUsecaseInApplication(appFile string) (map[string]int, error) {

	existingUsecaseNames := map[string]int{}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, appFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range astFile.Decls {

		ast.Inspect(decl, func(n ast.Node) bool {

			switch x := n.(type) {
			case *ast.CallExpr:
				z, ok := x.Fun.(*ast.SelectorExpr)
				if ok {
					if z.Sel.Name == "NewUsecase" {
						usecaseName, ok2 := z.X.(*ast.Ident)
						if ok2 {
							existingUsecaseNames[usecaseName.String()] = 1
						}
					}
				}
			}

			return true
		})
	}

	return existingUsecaseNames, nil

}

func findAllUsecaseInportNameFromController(domainName, controllerName string) ([]string, error) {

	res := make([]string, 0)

	folderName := fmt.Sprintf("%s/%s", getControllerFolder(domainName), strings.ToLower(controllerName))

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, folderName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	const HandlerSuffix = "handler"

	for _, pkg := range pkgs {

		// read file by file
		for _, astFile := range pkg.Files {

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
						res = append(res, uc)
					}
					return true
				})
			}

			//// in every declaration like type, func, const
			//for _, decl := range file.Decls {
			//
			//	// focus only to type
			//	gen, ok := decl.(*ast.GenDecl)
			//	if !ok || gen.Tok != token.TYPE {
			//		continue
			//	}
			//
			//	for _, specs := range gen.Specs {
			//
			//		ts, ok := specs.(*ast.TypeSpec)
			//		if !ok {
			//			continue
			//		}
			//
			//		if x, ok := ts.Type.(*ast.StructType); ok {
			//
			//			// check the specific struct name
			//			if !strings.HasSuffix(strings.ToLower(ts.Name.String()), "controller") {
			//				continue
			//			}
			//
			//			for _, field := range x.Fields.List {
			//
			//				se, ok := field.Type.(*ast.SelectorExpr)
			//				if !ok {
			//					continue
			//				}
			//
			//				if se.Sel.String() == "Inport" && len(field.Names) > 0 {
			//					name := field.Names[0].String()
			//					i := strings.Index(name, "Inport")
			//					res = append(res, name[:i])
			//
			//				}
			//
			//			}
			//
			//			//inportLine = fset.Position(iStruct.Fields.Closing).Line
			//			//return inportLine, nil
			//		}
			//	}
			//
			//}

		}

	}

	return res, nil
}

func InjectApplicationInMain(appName string) ([]byte, error) {

	templateWithData := fmt.Sprintf("\"%s\":application.New%s(),", strings.ToLower(appName), utils.PascalCase(appName))

	beforeLine, err := getInjectedLineInMain(appName)
	if err != nil {
		return nil, err
	}

	if beforeLine == 0 {
		return nil, nil
	}

	file, err := os.Open("main.go")
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

		if line == beforeLine-1 {
			buffer.WriteString(templateWithData)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
		line++
	}
	return buffer.Bytes(), nil
}

func getInjectedLineInMain(appName string) (int, error) {

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "main.go", nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}
	injectedLine := 0

	foundGogenRunner := false

	ast.Inspect(astFile, func(node ast.Node) bool {

		for {

			if foundGogenRunner {

				keyValueExpr, ok := node.(*ast.KeyValueExpr)
				if !ok {
					break
				}

				if keyValueExpr.Key.(*ast.BasicLit).Value == fmt.Sprintf("\"%s\"", appName) {
					injectedLine = 0
					return false
				}

			}

			compositLit, ok := node.(*ast.CompositeLit)
			if !ok {
				break
			}

			mapType, ok := compositLit.Type.(*ast.MapType)
			if !ok {
				break
			}

			selectorExpr, ok := mapType.Value.(*ast.SelectorExpr)
			if !ok {
				break
			}

			if selectorExpr.X.(*ast.Ident).String() == "gogen" && selectorExpr.Sel.String() == "Runner" {

				foundGogenRunner = true

				injectedLine = fset.Position(compositLit.Rbrace).Line

			}

			break
		}

		return true
	})

	return injectedLine, nil

}

func InjectRegisterUsecaseInApplication(usecaseName, appFile string) ([]byte, error) {

	templateWithData := fmt.Sprintf("%s.NewUsecase(datasource),", usecaseName)

	// u.AddUsecase(runordercreate.NewUsecase(datasource))
	//templateWithData := fmt.Sprintf("u.AddUsecase(%s.NewUsecase(datasource))", usecaseName)

	beforeLine, err := getInjectedLineInApplication(appFile)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(appFile)
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

		if line == beforeLine+1 {
			buffer.WriteString(templateWithData)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
		line++
	}
	return buffer.Bytes(), nil
}

func getInjectedLineInApplication(appFile string) (int, error) {

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, appFile, nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}
	injectedLine := 0

	ast.Inspect(astFile, func(node ast.Node) bool {

		//ast.Print(fset, node)

		for {

			exprStms, ok := node.(*ast.ExprStmt)
			if !ok {
				break
			}

			callExpr, ok := exprStms.X.(*ast.CallExpr)
			if !ok {
				break
			}

			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				break
			}

			if selectorExpr.Sel.String() == "AddUsecase" {
				injectedLine = fset.Position(selectorExpr.Sel.NamePos).Line
			}

			break
		}

		return true
	})

	return injectedLine, nil

}

func getControllerFolder(domainName string) string {
	return fmt.Sprintf("domain_%s/controller", domainName)
}

func getGatewayFolder(domainName string) string {
	return fmt.Sprintf("domain_%s/gateway", domainName)
}
