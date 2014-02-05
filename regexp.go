package data

import (
	"regexp"
)

var UserRegexp *regexp.Regexp
var IdentRegexp *regexp.Regexp
var PathRegexp *regexp.Regexp
var EmailRegexp *regexp.Regexp
var HandleRegexp *regexp.Regexp

func init() {
	identRE := "[A-Za-z0-9-_.]+"
	pathRE := "((" + identRE + ")/(" + identRE + "))"
	handleRE := pathRE + "(\\." + identRE + ")?(@" + identRE + ")?"
	emailRE := `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`

	UserRegexp = compileRegexp("^" + identRE + "$")
	IdentRegexp = compileRegexp("^" + identRE + "$")
	PathRegexp = compileRegexp("^" + pathRE + "$")
	EmailRegexp = compileRegexp("^" + emailRE + "$")
	HandleRegexp = compileRegexp("^" + handleRE + "$")
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
