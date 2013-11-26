package data

import (
	"fmt"
	"os"
)

var DEBUG bool

func DErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func DOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}
