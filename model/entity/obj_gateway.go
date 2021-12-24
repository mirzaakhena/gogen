package entity

import (
	"fmt"
	"github.com/mirzaakhena/gogen/infrastructure/util"
	"github.com/mirzaakhena/gogen/model/domerror"
	"github.com/mirzaakhena/gogen/model/vo"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const gatewayStructName = "gateway"

// ObjGateway  depend on (which) usecase that want to be tested
type ObjGateway struct {
	GatewayName vo.Naming
}

// ObjDataGateway  ...
type ObjDataGateway struct {
	PackagePath string
	GatewayName string
	Methods     vo.OutportMethods
}

// NewObjGateway   ...
func NewObjGateway(gatewayName string) (*ObjGateway, error) {

	if gatewayName == "" {
		return nil, domerror.GatewayNameMustNotEmpty
	}

	var obj ObjGateway
	obj.GatewayName = vo.Naming(gatewayName)

	return &obj, nil
}

// GetData ...
func (o ObjGateway) GetData(PackagePath string, outportMethods vo.OutportMethods) *ObjDataGateway {
	return &ObjDataGateway{
		PackagePath: PackagePath,
		GatewayName: o.GatewayName.LowerCase(),
		Methods:     outportMethods,
	}
}

// GetGatewayRootFolderName ...
func (o ObjGateway) GetGatewayRootFolderName() string {
	return fmt.Sprintf("gateway/%s", o.GatewayName.LowerCase())
}

// GetGatewayFileName ...
func (o ObjGateway) GetGatewayFileName() string {
	return fmt.Sprintf("%s/gateway.go", o.GetGatewayRootFolderName())
}

// GetGatewayStructName ...
func GetGatewayStructName() string {
	return fmt.Sprintf("gateway")
}

func (o ObjGateway) InjectToGateway(injectedCode string) ([]byte, error) {
	return InjectCodeAtTheEndOfFile(o.GetGatewayFileName(), injectedCode)
}

func FindGatewayByName(gatewayName string) (*ObjGateway, error) {
	folderName := fmt.Sprintf("gateway/%s", strings.ToLower(gatewayName))

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
						if !strings.HasSuffix(strings.ToLower(ts.Name.String()), gatewayStructName) {
							continue
						}

						return NewObjGateway(pkg.Name)

						//inportLine = fset.Position(iStruct.Fields.Closing).Line
						//return inportLine, nil
					}
				}

			}

		}

	}

	return nil, nil
}

func FindAllObjGateway() ([]*ObjGateway, error) {

	if !util.IsFileExist("gateway") {
		return nil, fmt.Errorf("gateway is not created yet")
	}

	dir, err := os.ReadDir("gateway")
	if err != nil {
		return nil, err
	}

	gateways := make([]*ObjGateway, 0)

	for _, d := range dir {
		if !d.IsDir() {
			continue
		}

		g, err := FindGatewayByName(d.Name())
		if err != nil {
			return nil, err
		}
		if g == nil {
			continue
		}

		gateways = append(gateways, g)

	}

	return gateways, nil
}
