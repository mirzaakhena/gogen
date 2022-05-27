package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// CamelCase is
func CamelCase(name string) string {

	// TODO
	// hardcoded is bad
	// But we can figure out later
	{
		if name == "IPAddress" {
			return "ipAddress"
		}

		if name == "ID" {
			return "id"
		}
	}

	out := []rune(name)
	out[0] = unicode.ToLower([]rune(name)[0])
	return string(out)
}

// UpperCase is
func UpperCase(name string) string {
	return strings.ToUpper(name)
}

// LowerCase is
func LowerCase(name string) string {
	return strings.ToLower(name)
}

var matchFirstCapSpaceCase = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCapSpaceCase = regexp.MustCompile("([a-z0-9])([A-Z])")

// SpaceCase is
func SpaceCase(str string) string {
	snake := matchFirstCapSpaceCase.ReplaceAllString(str, "${1} ${2}")
	snake = matchAllCapSpaceCase.ReplaceAllString(snake, "${1} ${2}")
	return strings.ToLower(snake)
}

// PascalCase is
func PascalCase(name string) string {
	rs := []rune(name)
	return strings.ToUpper(string(rs[0])) + string(rs[1:])
}

var matchFirstCapSnakeCase = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCapSnakeCase = regexp.MustCompile("([a-z0-9])([A-Z])")

// SnakeCase is
func SnakeCase(str string) string {
	snake := matchFirstCapSnakeCase.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCapSnakeCase.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// StartWith is
func StartWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// HasOneOfThisPrefix ...
func HasOneOfThisPrefix(s string, str ...string) bool {
	lc := LowerCase(s)
	for _, s := range str {
		if strings.HasPrefix(lc, s) {
			return true
		}
	}
	return false
}
