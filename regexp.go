package data

import (
	"regexp"
)

var UserRegexp *regexp.Regexp
var NameRegexp *regexp.Regexp
var PathRegexp *regexp.Regexp
var HandleRegexp *regexp.Regexp

func init() {
	identRE := "[A-Za-z0-9-_.]+"
	pathRE := "((" + identRE + ")/(" + identRE + "))"
	handleRE := "^" + pathRE + "(\\." + identRE + ")?(@" + identRE + ")?$"

	UserRegexp = compileRegexp(identRE)
	NameRegexp = compileRegexp(identRE)
	PathRegexp = compileRegexp(pathRE)
	HandleRegexp = compileRegexp(handleRE)
}

func compileRegexp(s string) *regexp.Regexp {
	r, err := regexp.Compile(s)
	if err != nil {
		pOut("%s", err)
		pOut("%v", r)
		panic("Regex does not compile: " + s)
	}
	return r
}
