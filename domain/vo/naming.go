package vo

import (
  "github.com/mirzaakhena/gogen/infrastructure/util"
	"strings"
)

// Naming ...
type Naming string

// String ...
func (r Naming) String() string {
	return string(r)
}

// HasPrefix ...
func (r Naming) HasPrefix(str string) bool {
	return strings.HasPrefix(strings.ToLower(string(r)), str)
}

// HasOneOfThisPrefix ...
func (r Naming) HasOneOfThisPrefix(str ...string) bool {
	lc := r.LowerCase()
	for _, s := range str {
		if strings.HasPrefix(lc, s) {
			return true
		}
	}
	return false
}

// IsEmpty ...
func (r Naming) IsEmpty() bool {
	return len(string(r)) == 0
}

// CamelCase is
func (r Naming) CamelCase() string {
	return util.CamelCase(string(r))
}

// UpperCase is
func (r Naming) UpperCase() string {
	return util.UpperCase(string(r))
}

// LowerCase is
func (r Naming) LowerCase() string {
	return util.LowerCase(string(r))
}

// SpaceCase is
func (r Naming) SpaceCase() string {
	return util.SpaceCase(string(r))
}

// PascalCase is
func (r Naming) PascalCase() string {
	return util.PascalCase(string(r))
}

// SnakeCase is
func (r Naming) SnakeCase() string {
	return util.SnakeCase(string(r))
}
