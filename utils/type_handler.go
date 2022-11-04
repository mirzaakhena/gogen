package utils

import (
	"bytes"
	"fmt"
	"go/ast"
)

// TypeHandler ...
type TypeHandler struct {
	PrefixExpression string //
}

func (r TypeHandler) Start(expr ast.Expr) string {
	return r.appendType(expr)
}

func (r TypeHandler) appendType(expr ast.Expr) string {
	var param bytes.Buffer

	for {
		switch t := expr.(type) {
		case *ast.Ident:
			return r.processIdent(&param, t)

		case *ast.ArrayType:
			return r.processArrayType(&param, t)

		case *ast.StarExpr:
			return r.processStarExpr(&param, t)

		case *ast.SelectorExpr:
			return r.processSelectorExpr(&param, t)

		case *ast.InterfaceType:
			return r.processInterfaceType(&param, t)

		case *ast.ChanType:
			return r.processChanType(&param, t)

		case *ast.MapType:
			return r.processMapType(&param, t)

		case *ast.FuncType:
			param.WriteString("func")
			return r.processFuncType(&param, t)

		case *ast.StructType:
			return r.processStruct(&param, t)

		default:
			return param.String()
		}

	}
}

func (r TypeHandler) processIdent(param *bytes.Buffer, t *ast.Ident) string {
	if t.Obj != nil {
		param.WriteString(r.PrefixExpression)
		param.WriteString(".")
	}
	param.WriteString(t.Name)
	return param.String()
}

func (r TypeHandler) processStruct(param *bytes.Buffer, t *ast.StructType) string {
	param.WriteString("struct{")

	nParam := t.Fields.NumFields()
	c := 0
	for _, field := range t.Fields.List {
		nNames := len(field.Names)
		c += nNames
		for i, name := range field.Names {
			param.WriteString(name.String())
			if i < len(field.Names)-1 {
				param.WriteString(", ")
			} else {
				param.WriteString(" ")
			}
		}
		param.WriteString(r.appendType(field.Type))

		if c < nParam {
			param.WriteString("; ")
		}

	}

	param.WriteString("}")
	return param.String()
}

func (r TypeHandler) processArrayType(param *bytes.Buffer, t *ast.ArrayType) string {
	if t.Len != nil {
		arrayCapacity := t.Len.(*ast.BasicLit).Value
		param.WriteString(fmt.Sprintf("[%s]", arrayCapacity))
	} else {
		param.WriteString("[]")
	}
	param.WriteString(r.appendType(t.Elt))

	return param.String()
}

func (r TypeHandler) processStarExpr(param *bytes.Buffer, t *ast.StarExpr) string {
	param.WriteString("*")
	param.WriteString(r.appendType(t.X))
	return param.String()
}

func (r TypeHandler) processInterfaceType(param *bytes.Buffer, t *ast.InterfaceType) string {
	param.WriteString("any")
	_ = t
	return param.String()
}

func (r TypeHandler) processSelectorExpr(param *bytes.Buffer, t *ast.SelectorExpr) string {
	param.WriteString(r.appendType(t.X))
	param.WriteString(".")
	param.WriteString(t.Sel.Name)
	return param.String()
}

func (r TypeHandler) processChanType(param *bytes.Buffer, t *ast.ChanType) string {
	if t.Dir == 1 {
		param.WriteString("chan<- ")
	} else if t.Dir == 2 {
		param.WriteString("<-chan ")
	} else {
		param.WriteString("chan ")
	}

	param.WriteString(r.appendType(t.Value))
	return param.String()
}

func (r TypeHandler) processMapType(param *bytes.Buffer, t *ast.MapType) string {
	param.WriteString("map[")
	param.WriteString(r.appendType(t.Key))
	param.WriteString("]")
	param.WriteString(r.appendType(t.Value))
	return param.String()
}

func (r TypeHandler) processFuncType(param *bytes.Buffer, t *ast.FuncType) string {

	// TODO need to handle method param without variable
	// TODO need to handle param with struct/interface type

	nParam := t.Params.NumFields()
	param.WriteString("(")
	for iList, field := range t.Params.List {
		nNames := len(field.Names)
		for i, name := range field.Names {
			param.WriteString(name.String())
			if i < len(field.Names)-1 {
				param.WriteString(", ")
			} else {
				param.WriteString(" ")
			}
		}
		//param.WriteString(r.PrefixExpression)
		//param.WriteString(".")
		param.WriteString(r.appendType(field.Type))
		if iList+nNames < nParam {
			param.WriteString(", ")
		}

	}
	param.WriteString(") ")

	nResult := t.Results.NumFields()

	haveParentThesis := false
	if t.Results == nil {
		return param.String()
	}
	for i, field := range t.Results.List {

		if i == 0 {
			if len(field.Names) > 0 || nResult > 1 {
				param.WriteString("(")
				haveParentThesis = true
			}
		}

		for i, name := range field.Names {
			param.WriteString(name.String() + " ")
			if i < len(field.Names)-1 {
				param.WriteString(", ")
			}
		}
		param.WriteString(r.appendType(field.Type))
		if i+1 < nResult {
			param.WriteString(", ")
		}
	}

	if haveParentThesis {
		param.WriteString(")")
	}

	return param.String()
}
