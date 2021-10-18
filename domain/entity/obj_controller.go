package entity

import (
  "bufio"
  "bytes"
  "fmt"
  "github.com/mirzaakhena/gogen/application/apperror"
  "github.com/mirzaakhena/gogen/domain/vo"
  "github.com/mirzaakhena/gogen/infrastructure/util"
  "go/ast"
  "go/parser"
  "go/token"
  "os"
  "strings"
)

type ObjController struct {
  ControllerName vo.Naming
  DriverName     string
}

// ObjDataController ...
type ObjDataController struct {
  PackagePath    string
  UsecaseName    string
  ControllerName string
}

func NewObjController(controllerName, driverName string) (*ObjController, error) {

  if controllerName == "" {
    return nil, apperror.ControllerNameMustNotEmpty
  }

  var obj ObjController
  obj.ControllerName = vo.Naming(controllerName)
  obj.DriverName = driverName

  return &obj, nil
}

// GetData ...
func (o ObjController) GetData(PackagePath string, ou ObjUsecase) *ObjDataController {
  return &ObjDataController{
    PackagePath:    PackagePath,
    ControllerName: o.ControllerName.String(),
    UsecaseName:    ou.UsecaseName.String(),
  }
}

// GetControllerRootFolderName ...
func (o ObjController) GetControllerRootFolderName() string {
  return fmt.Sprintf("controller/%s", o.ControllerName.LowerCase())
}

// GetControllerInterfaceFile ...
func (o ObjController) GetControllerInterfaceFile() string {
  return fmt.Sprintf("controller/controller.go")
}

// GetControllerResponseFileName ...
func (o ObjController) GetControllerResponseFileName() string {
  return fmt.Sprintf("%s/response.go", o.GetControllerRootFolderName())
}

// GetControllerInterceptorFileName ...
func (o ObjController) GetControllerInterceptorFileName() string {
  return fmt.Sprintf("%s/interceptor.go", o.GetControllerRootFolderName())
}

// GetControllerRouterFileName ...
func (o ObjController) GetControllerRouterFileName() string {
  return fmt.Sprintf("%s/router.go", o.GetControllerRootFolderName())
}

// GetControllerHandlerFileName ...
func (o ObjController) GetControllerHandlerFileName(ou ObjUsecase) string {
  return fmt.Sprintf("%s/handler_%s.go", o.GetControllerRootFolderName(), ou.UsecaseName.LowerCase())
}

func (o ObjController) InjectInportToStruct(templateWithData string) ([]byte, error) {

  inportLine, err := o.getInportLine()
  if err != nil {
    return nil, err
  }

  controllerFile := o.GetControllerRouterFileName()

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

func (o ObjController) InjectRouterBind(templateWithData string) ([]byte, error) {

  controllerFile := o.GetControllerRouterFileName()

  routerLine, err := o.getBindRouterLine()
  if err != nil {
    return nil, err
  }

  //templateCode, err := util.PrintTemplate(templates.ControllerBindRouterGinFile, obj)
  //if err != nil {
  //  return err
  //}

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

func (o ObjController) getInportLine() (int, error) {

  controllerFile := o.GetControllerRouterFileName()

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

func (o ObjController) getBindRouterLine() (int, error) {

  controllerFile := o.GetControllerRouterFileName()

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

const controllerStructName = "controller"

func FindControllerByName(controllerName string) (*ObjController, error) {
  folderName := fmt.Sprintf("controller/%s", strings.ToLower(controllerName))

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

          if _, ok := ts.Type.(*ast.StructType); ok {

            // check the specific struct name
            if !strings.HasSuffix(strings.ToLower(ts.Name.String()), controllerStructName) {
              continue
            }

            return NewObjController(pkg.Name, "")

            //inportLine = fset.Position(iStruct.Fields.Closing).Line
            //return inportLine, nil
          }
        }

      }

    }

  }

  return nil, nil
}

func FindAllObjController() ([]*ObjController, error) {

  if !util.IsFileExist("controller") {
    return nil, fmt.Errorf("controller is not created yet")
  }

  dir, err := os.ReadDir("controller")
  if err != nil {
    return nil, err
  }

  controllers := make([]*ObjController, 0)

  for _, d := range dir {
    if !d.IsDir() {
      continue
    }

    g, err := FindControllerByName(d.Name())
    if err != nil {
      return nil, err
    }
    if g == nil {
      continue
    }

    controllers = append(controllers, g)

  }

  return controllers, nil
}

func (o ObjController) FindAllUsecaseInportNameFromController() ([]string, error) {

  res := make([]string, 0)

  folderName := fmt.Sprintf("controller/%s", strings.ToLower(o.ControllerName.LowerCase()))

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
            if !strings.HasSuffix(strings.ToLower(ts.Name.String()), controllerStructName) {
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
