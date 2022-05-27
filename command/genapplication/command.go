package genapplication

import (
	"fmt"
	"github.com/mirzaakhena/gogen/utils"
	"go/ast"
	"go/parser"
	"go/token"

	"os"
	"strings"
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

	domainName := utils.GetDefaultDomain()
	applicationName := inputs[0]

	obj := &ObjTemplate{
		PackagePath:     utils.GetPackagePath(),
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

		usecaseNames, err := findAllUsecaseInportNameFromController(domainName, objController)
		if err != nil {
			return err
		}

		obj.GatewayName = &objGateway
		obj.ControllerName = &objController
		obj.UsecaseNames = usecaseNames

		fileRenamer := map[string]string{
			"domainname":      utils.LowerCase(domainName),
			"applicationname": utils.LowerCase(applicationName),
		}

		err = utils.CreateEverythingExactly("templates/application/", driverName, fileRenamer, obj, utils.AppTemplates)
		if err != nil {
			return err
		}

	} //---------

	// inject to main.__go
	{
		fset := token.NewFileSet()
		utils.InjectToMain(fset, applicationName)
	}

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

func findAllUsecaseInportNameFromController(domainName, controllerName string) ([]string, error) {

	res := make([]string, 0)

	folderName := fmt.Sprintf("%s/%s", getControllerFolder(domainName), strings.ToLower(controllerName))

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, folderName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
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

					if x, ok := ts.Type.(*ast.StructType); ok {

						// check the specific struct name
						if !strings.HasSuffix(strings.ToLower(ts.Name.String()), "controller") {
							continue
						}

						for _, field := range x.Fields.List {

							se, ok := field.Type.(*ast.SelectorExpr)
							if !ok {
								continue
							}

							if se.Sel.String() == "Inport" && len(field.Names) > 0 {
								name := field.Names[0].String()
								i := strings.Index(name, "Inport")
								res = append(res, name[:i])

							}

						}

						//inportLine = fset.Position(iStruct.Fields.Closing).Line
						//return inportLine, nil
					}
				}

			}

		}

	}

	return res, nil
}

func getControllerFolder(domainName string) string {
	return fmt.Sprintf("domain_%s/controller", domainName)
}

func getGatewayFolder(domainName string) string {
	return fmt.Sprintf("domain_%s/gateway", domainName)
}
