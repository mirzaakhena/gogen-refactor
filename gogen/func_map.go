package gogen

import (
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

var FuncMap = template.FuncMap{
	"camelCase":  camelCase,
	"PascalCase": PascalCase,
	"SnakeCase":  snakeCase,
	"UpperCase":  upperCase,
	"LowerCase":  lowerCase,
	"SpaceCase":  SpaceCase,
	"StartWith":  startWith,
}

func camelCase(name string) string {

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

func upperCase(name string) string {
	return strings.ToUpper(name)
}

func lowerCase(name string) string {
	return strings.ToLower(name)
}

var matchFirstCapSpaceCase = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCapSpaceCase = regexp.MustCompile("([a-z0-9])([A-Z])")

func SpaceCase(str string) string {
	snake := matchFirstCapSpaceCase.ReplaceAllString(str, "${1} ${2}")
	snake = matchAllCapSpaceCase.ReplaceAllString(snake, "${1} ${2}")
	return strings.ToLower(snake)
}

func PascalCase(name string) string {
	rs := []rune(name)
	return strings.ToUpper(string(rs[0])) + string(rs[1:])
}

var matchFirstCapSnakeCase = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCapSnakeCase = regexp.MustCompile("([a-z0-9])([A-Z])")

func snakeCase(str string) string {
	snake := matchFirstCapSnakeCase.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCapSnakeCase.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func startWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}
