package data

import (
	"os"
)

func helpCmd([]string) error {
	Usage()
	return nil
}

func Usage() {
	pErr(usageStr1)
	PrintCommands()
	pErr("\n")
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
