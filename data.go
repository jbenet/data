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
		fmt.Println("debugging on")
	}

	RegisterCommands()

	if len(os.Args) < 2 {
		Usage()
	}

	DispatchCommand(os.Args[1])
}
