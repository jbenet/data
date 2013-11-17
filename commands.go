package main

import (
	"fmt"
	"os"
)

type Command struct {
	name    string
	desc    string
	handler func()
}

var commands = map[string]Command{}

func RegisterCommand(name string, desc string, handler func()) {
	commands[name] = Command{name, desc, handler}
}

func RegisterCommands() {
	RC := RegisterCommand
	RC("version", "Show data version information.", VersionCmd)
	RC("help", "Show usage information.", HelpCmd)
}

func PrintCommands() {
	for _, cmd := range commands {
		fmt.Fprintf(os.Stderr, "    %-10.10s%s\n", cmd.name, cmd.desc)
	}
}

func DispatchCommand(name string) {

	if DEBUG {
		fmt.Fprintf(os.Stderr, "dispatching command %s\n", name)
	}

	cmd, ok := commands[name]
	if ok {
		cmd.handler()
	} else {
		fmt.Fprintf(os.Stderr, "data: unknown command \"%s\"\n", name)
		fmt.Fprintf(os.Stderr, "Run `data help` for usage.\n")
	}
}
