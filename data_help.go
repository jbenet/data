package data

import (
	"os"
)

func HelpCmd([]string) error {
	Usage()
	return nil
}

func Usage() {
	DErr(usageStr1)
	PrintCommands()
	DErr("\n")
	// DErr(usageStr2)

	os.Exit(1)
}

var usageStr1 = `data is a dataset package manager.

Usage:

    data <command> [arguments]

Commands:

`

var usageStr2 = `Use 'data help <command>' for more information about a command.

`
