package prod

import (
  "github.com/mirzaakhena/gogen/infrastructure/util"
	"text/template"
)

//FuncMap that used in templates
var FuncMap = template.FuncMap{
	"CamelCase":  util.CamelCase,
	"PascalCase": util.PascalCase,
	"SnakeCase":  util.SnakeCase,
	"UpperCase":  util.UpperCase,
	"LowerCase":  util.LowerCase,
	"SpaceCase":  util.SpaceCase,
}
