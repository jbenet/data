package main

import (
	"fmt"
	"os"
)

func HelpCmd() {
	Usage()
}

func Usage() {
	fmt.Fprintf(os.Stderr, usageStr1)
	PrintCommands()
	fmt.Fprintf(os.Stderr, "\n")
	// fmt.Fprintf(os.Stderr, usageStr2)

	os.Exit(1)
}

var usageStr1 = `data is a dataset package manager.

Usage:

    data <command> [arguments]

Commands:

`

var usageStr2 = `Use 'data help <command>' for more information about a command.

`
