package main

import (
	"flag"
	"fmt"
	"os"
)

var DEBUG bool

func main() {

	flag.BoolVar(&DEBUG, "debug", false, "Debug mode")
	flag.Parse()

	if DEBUG {
		DOut("debugging on\n")
	}

	RegisterCommands()

	if len(os.Args) < 2 {
		Usage()
	}

	DispatchCommand(os.Args[1])
}

func DErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func DOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}
