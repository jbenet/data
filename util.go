package data

import (
	"fmt"
	"os"
)

var DEBUG bool

func Err(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func Out(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func DErr(format string, a ...interface{}) {
	if DEBUG {
		Err(format, a...)
	}
}

func DOut(format string, a ...interface{}) {
	if DEBUG {
		Out(format, a...)
	}
}
