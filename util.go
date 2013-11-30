package data

import (
	"fmt"
	"os"
)

var DEBUG bool

func pErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func pOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func dErr(format string, a ...interface{}) {
	if DEBUG {
		pErr(format, a...)
	}
}

func dOut(format string, a ...interface{}) {
	if DEBUG {
		pOut(format, a...)
	}
}
