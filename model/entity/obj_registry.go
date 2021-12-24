package entity

import (
	"fmt"
	"github.com/mirzaakhena/gogen/model/vo"
	"go/ast"
	"go/parser"
	"go/token"
)

type ObjRegistry struct {
	RegistryName  vo.Naming
	ObjController *ObjController
	ObjGateway    *ObjGateway
	UsecaseNames  []string
}

type ObjDataRegistry struct {
	PackagePath    string
	RegistryName   string
	ControllerName string
	UsecaseNames   []string
	GatewayName    string
}

type ObjGatewayRequest struct {
	RegistryName  string
	ObjController *ObjController
	ObjGateway    *ObjGateway
	UsecaseNames  []string
}

func NewObjRegistry(req ObjGatewayRequest) (*ObjRegistry, error) {

	if req.ObjController == nil {
		return nil, fmt.Errorf("ObjController must not be nil")
	}

	if req.ObjGateway == nil {
		return nil, fmt.Errorf("ObjGateway must not be nil")
	}

	if len(req.UsecaseNames) == 0 {
		return nil, fmt.Errorf("usecases must not empty")
	}

	var obj ObjRegistry
	obj.RegistryName = vo.Naming(req.RegistryName)
	obj.ObjController = req.ObjController
	obj.ObjGateway = req.ObjGateway
	obj.UsecaseNames = req.UsecaseNames

	return &obj, nil
}

// GetData ...
func (o ObjRegistry) GetData(PackagePath string) *ObjDataRegistry {

	return &ObjDataRegistry{
		PackagePath:    PackagePath,
		RegistryName:   o.RegistryName.String(),
		ControllerName: o.ObjController.ControllerName.String(),
		GatewayName:    o.ObjGateway.GatewayName.LowerCase(),
		UsecaseNames:   o.UsecaseNames,
	}
}

// GetRegistryRootFolderName ...
func GetRegistryRootFolderName() string {
	return fmt.Sprintf("application/registry")
}

// GetApplicationFileName ...
func GetApplicationFileName() string {
	return fmt.Sprintf("application/application.go")
}

// GetRegistryFileName ...
func (o ObjRegistry) GetRegistryFileName() string {
	return fmt.Sprintf("%s/%s.go", GetRegistryRootFolderName(), o.RegistryName.LowerCase())
}

func (o ObjRegistry) InjectUsecaseInportField() error {

	fset := token.NewFileSet()

	//pkgs, err := parser.ParseDir(fset, o.GetRegistryFileName(), nil, parser.ParseComments)
	//if err != nil {
	//  return err
	//}

	file, err := parser.ParseFile(fset, o.GetRegistryFileName(), nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// in every declaration like type, func, const
	for _, decl := range file.Decls {

		// focus only to type
		//fgen, ok := decl.(*ast.FuncDecl)
		//if !ok || gen.Tok != token.TYPE {
		// continue
		//}

		ast.Print(fset, decl)

		//for _, specs := range gen.Specs {
		//
		//  ts, ok := specs.(*ast.TypeSpec)
		//  if !ok {
		//    continue
		//  }
		//
		//  if _, ok := ts.Type.(*ast.StructType); ok {
		//
		//    // check the specific struct name
		//    if !strings.HasSuffix(strings.ToLower(ts.Name.String()), gatewayStructName) {
		//      continue
		//    }
		//
		//    return NewObjGateway(pkg.Name)
		//
		//    //inportLine = fset.Position(iStruct.Fields.Closing).Line
		//    //return inportLine, nil
		//  }
		//}

	}

	//var iFace *ast.InterfaceType
	//var file *ast.File

	//isAlreadyInjectedBefore := IsExist(fset, o.ObjUsecase.GetUsecaseRootFolderName(), func(f *ast.File, ts *ast.TypeSpec) bool {
	//  interfaceType, ok := ts.Type.(*ast.InterfaceType)
	//
	//  if ok && ts.Name.String() == OutportInterfaceName {
	//
	//    file = f
	//    iFace = interfaceType
	//
	//    // trace every method to find something line like `repository.SaveOrderRepo`
	//    for _, meth := range iFace.Methods.List {
	//
	//      selType, ok := meth.Type.(*ast.SelectorExpr)
	//
	//      // if interface already injected (SaveOrderRepo) then abort the mission
	//      if ok && selType.Sel.String() == o.getRepositoryName() {
	//        return true
	//      }
	//    }
	//  }
	//  return false
	//})
	//
	//// if it is already injected, then nothing to do
	//if isAlreadyInjectedBefore {
	//  return nil
	//}
	//
	//if iFace == nil {
	//  return fmt.Errorf("outport struct not found")
	//}
	//
	//// we want to inject it now
	//// add new repository to outport interface
	//iFace.Methods.List = append(iFace.Methods.List, &ast.Field{
	//  Type: &ast.SelectorExpr{
	//    X: &ast.Ident{
	//      Name: GetPackageName(o.GetRepositoryRootFolderName()),
	//    },
	//    Sel: &ast.Ident{
	//      Name: o.getRepositoryName(),
	//    },
	//  },
	//})
	//
	//// TODO who is responsible to write a file? entity or gateway?
	//// i prefer to use gateway instead of entity
	//
	//fileReadPath := fset.Position(file.Pos()).Filename
	//
	//// rewrite the outport
	//f, err := os.Create(fileReadPath)
	//if err != nil {
	//  return err
	//}
	//
	//if err := printer.Fprint(f, fset, file); err != nil {
	//  return err
	//}
	//err = f.Close()
	//if err != nil {
	//  return err
	//}
	//
	//// reformat and import
	//newBytes, err := imports.Process(fileReadPath, nil, nil)
	//if err != nil {
	//  return err
	//}
	//
	//if err := ioutil.WriteFile(fileReadPath, newBytes, 0644); err != nil {
	//  return err
	//}

	return nil
}
