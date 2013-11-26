package main

import (
	"github.com/jbenet/data"
	"flag"
	"os"
)

func main() {

	flag.BoolVar(&data.DEBUG, "debug", false, "Debug mode")
	flag.Parse()

	if data.DEBUG {
		data.DOut("debugging on\n")
	}

	data.RegisterCommands()

	if len(os.Args) < 2 {
		data.Usage()
	}

	data.DispatchCommand(os.Args[1], os.Args[2:])
}
