package data

import (
	"os"
)

func HelpCmd([]string) error {
	Usage()
	return nil
}

func Usage() {
	Err(usageStr1)
	PrintCommands()
	Err("\n")
	// Err(usageStr2)

	os.Exit(1)
}

var usageStr1 = `data is a dataset package manager.

Usage:

    data <command> [arguments]

Commands:

`

var usageStr2 = `Use 'data help <command>' for more information about a command.

`
