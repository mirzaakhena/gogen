package utils

import "text/template"

var FuncMap = template.FuncMap{
	"CamelCase":  CamelCase,
	"PascalCase": PascalCase,
	"SnakeCase":  SnakeCase,
	"UpperCase":  UpperCase,
	"LowerCase":  LowerCase,
	"SpaceCase":  SpaceCase,
	"StartWith":  StartWith,
}
