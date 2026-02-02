package main

import (
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var fns = template.FuncMap{
	"last": func(x int, a any) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
	"title": StringToTitle,
	"replace": func(s, old, new string) string {
		return strings.ReplaceAll(s, old, new)
	},
}

func StringToTitle(s string) string {
	c := cases.Title(language.English)
	return c.String(s)
}
