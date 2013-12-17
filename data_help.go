package data

import (
	"os"
	"strings"
)

func helpCmd(args []string) error {
	args = append([]string{"data"}, args...)
	group := strings.Join(args, " ")
	Usage(group)
	return nil
}

func Usage(group string) {
	pErr(usageStr1)
	PrintCommands(group)
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
