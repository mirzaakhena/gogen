package gogen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

func ReadImports(node *ast.File) map[string]int {
	importSpecs := map[string]int{}
	for index, importSpec := range node.Imports {
		importSpecs[importSpec.Path.Value] = index
	}
	return importSpecs
}

func appendWithOrder(ss []string, s string) []string {
	i := sort.SearchStrings(ss, s)
	ss = append(ss, "")
	copy(ss[i+1:], ss[i:])
	ss[i] = s
	return ss
}

func ReadAllStructInFile(node *ast.File, structureStructs map[string][]FieldType) (map[string][]FieldType, error) {

	for _, dec := range node.Decls {
		if gen, ok := dec.(*ast.GenDecl); ok {
			if gen.Tok != token.TYPE {
				continue
			}
			for _, specs := range gen.Specs {
				if ts, ok := specs.(*ast.TypeSpec); ok {
					if istruct, ok := ts.Type.(*ast.StructType); ok {

						fieldTypes := []FieldType{}
						for _, field := range istruct.Fields.List {
							tag := ""
							if field.Tag != nil {
								tag = field.Tag.Value
							}
							comment := ""
							if field.Comment != nil {
								comment = field.Comment.List[0].Text
							}
							for _, fieldName := range field.Names {
								fieldTypes = append(fieldTypes, FieldType{
									FieldName: fieldName.Name,
									Type:      appendType(field.Type),
									Tag:       tag,
									Comment:   comment,
								})
							}
						}

						if _, exist := structureStructs[ts.Name.String()]; exist {
							return nil, fmt.Errorf("struct %s is exist", ts.Name.String())
						}

						structureStructs[ts.Name.String()] = fieldTypes
					}
				}
			}
		}
	}

	return structureStructs, nil
}

func ReadFieldInStruct(node *ast.File, structName string) []FieldType {

	for _, dec := range node.Decls {
		if gen, ok := dec.(*ast.GenDecl); ok {
			if gen.Tok != token.TYPE {
				continue
			}
			for _, specs := range gen.Specs {
				if ts, ok := specs.(*ast.TypeSpec); ok {
					if istruct, ok := ts.Type.(*ast.StructType); ok {
						if !strings.HasSuffix(ts.Name.String(), structName) {
							continue
						}
						fieldTypes := []FieldType{}
						for _, field := range istruct.Fields.List {
							tag := ""
							if field.Tag != nil {
								tag = field.Tag.Value
							}
							comment := ""
							if field.Comment != nil {
								comment = field.Comment.List[0].Text
							}
							for _, fieldName := range field.Names {
								fieldTypes = append(fieldTypes, FieldType{
									FieldName: fieldName.Name,
									Type:      appendType(field.Type),
									Tag:       tag,
									Comment:   comment,
								})
							}

							// ast.Print(fset, field)
						}
						return fieldTypes
					}
				}
			}
		}
	}
	return nil
}

func ReadInterfaceMethodAndField(node *ast.File, interfaceName string, mapStruct map[string][]FieldType) ([]InterfaceMethod, error) {

	for _, dec := range node.Decls {
		if gen, ok := dec.(*ast.GenDecl); ok {
			if gen.Tok != token.TYPE {
				continue
			}
			for _, specs := range gen.Specs {
				if ts, ok := specs.(*ast.TypeSpec); ok {
					if iface, ok := ts.Type.(*ast.InterfaceType); ok {
						if ts.Name.String() != interfaceName {
							continue
						}
						methods := []InterfaceMethod{}
						for _, meths := range iface.Methods.List {

							fType := meths.Type.(*ast.FuncType)
							if len(fType.Params.List) < 2 {
								return nil, fmt.Errorf("Need second params")
							}

							if fType.Results.NumFields() < 1 {
								return nil, fmt.Errorf("Need result params")
							}

							paramType := fType.Params.List[1].Type.(*ast.Ident).Name
							resultType := fType.Results.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name

							methods = append(methods, InterfaceMethod{
								MethodName:     meths.Names[0].String(),
								ParamType:      paramType,
								ResultType:     resultType,
								RequestFields:  mapStruct[paramType],
								ResponseFields: mapStruct[resultType],
							})

						}
						return methods, nil

					}
				}
			}
		}
	}
	return nil, fmt.Errorf("interface %s not found", interfaceName)
}

func ReadInterfaceName(node *ast.File) ([]string, error) {

	iNames := []string{}

	for _, dec := range node.Decls {
		if gen, ok := dec.(*ast.GenDecl); ok {
			if gen.Tok != token.TYPE {
				continue
			}
			for _, specs := range gen.Specs {
				if ts, ok := specs.(*ast.TypeSpec); ok {
					if _, ok := ts.Type.(*ast.InterfaceType); ok {
						iNames = append(iNames, ts.Name.String())
					}
				}
			}
		}
	}
	return iNames, nil
}

func appendType(expr ast.Expr) string {
	var param bytes.Buffer

	for {
		switch t := expr.(type) {
		case *ast.Ident:
			return processIdent(&param, t)

		case *ast.ArrayType:
			return processArrayType(&param, t)

		case *ast.StarExpr:
			return processStarExpr(&param, t)

		case *ast.SelectorExpr:
			return processSelectorExpr(&param, t)

		case *ast.InterfaceType:
			return processInterfaceType(&param, t)

		case *ast.ChanType:
			return processChanType(&param, t)

		case *ast.MapType:
			return processMapType(&param, t)

		case *ast.FuncType:
			return processFuncType(&param, t)

		default:
			return param.String()
		}

	}
}

func processIdent(param *bytes.Buffer, t *ast.Ident) string {
	param.WriteString(t.Name)
	return param.String()
}

func processArrayType(param *bytes.Buffer, t *ast.ArrayType) string {
	if t.Len != nil {
		arrayCapacity := t.Len.(*ast.BasicLit).Value
		param.WriteString(fmt.Sprintf("[%s]", arrayCapacity))
	} else {
		param.WriteString("[]")
	}
	param.WriteString(appendType(t.Elt))

	return param.String()
}

func processStarExpr(param *bytes.Buffer, t *ast.StarExpr) string {
	param.WriteString("*")
	param.WriteString(appendType(t.X))
	return param.String()
}

func processInterfaceType(param *bytes.Buffer, t *ast.InterfaceType) string {
	param.WriteString("interface{}")
	return param.String()
}

func processSelectorExpr(param *bytes.Buffer, t *ast.SelectorExpr) string {
	param.WriteString(appendType(t.X))
	param.WriteString(".")
	param.WriteString(t.Sel.Name)
	return param.String()
}

func processChanType(param *bytes.Buffer, t *ast.ChanType) string {
	if t.Dir == 1 {
		param.WriteString("chan<- ")
	} else if t.Dir == 2 {
		param.WriteString("<-chan ")
	} else {
		param.WriteString("chan ")
	}

	param.WriteString(appendType(t.Value))
	return param.String()
}

func processMapType(param *bytes.Buffer, t *ast.MapType) string {
	param.WriteString("map[")
	param.WriteString(appendType(t.Key))
	param.WriteString("]")
	param.WriteString(appendType(t.Value))
	return param.String()
}

func processFuncType(param *bytes.Buffer, t *ast.FuncType) string {
	param.WriteString("func(")

	nParam := t.Params.NumFields()
	for i, field := range t.Params.List {
		param.WriteString(appendType(field.Type))
		if i+1 < nParam {
			param.WriteString(", ")
		} else {
			param.WriteString(") ")
		}
	}

	nResult := t.Results.NumFields()

	if nResult > 1 {
		param.WriteString("(")
	}

	for i, field := range t.Results.List {
		param.WriteString(appendType(field.Type))
		if i+1 < nResult {
			param.WriteString(", ")
		}
	}

	if nResult > 1 {
		param.WriteString(")")
	}

	return param.String()
}
