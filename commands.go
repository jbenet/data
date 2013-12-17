package data

import (
	"strings"
)

type CommandFunc func([]string) error

type Command struct {
	name    string
	desc    string
	handler CommandFunc
}

var commands = map[string]Command{}

func RegisterCommand(name string, desc string, handler CommandFunc) {
	commands[name] = Command{name, desc, handler}
}

func RegisterCommands() {
	RC := RegisterCommand
	RC("data", "Dataset package manager tool.", helpCmd)
	RC("data get", "Download and install dataset.", getCmd)
	RC("data list", "List installed datasets.", listCmd)
	RC("data info", "Show dataset information.", infoCmd)
	RC("data help", "Show usage information.", helpCmd)
	RC("data version", "Show data version information.", versionCmd)
	RC("data upload", "Upload dataset to storage service.", uploadCmd)
	RC("data manifest", "Generate dataset manifest.", manifestCmd)
}

func PrintCommands(group string) {

	if group != "data" {
		cmd := commands[group]
		pErr("%s: %s\n\n", cmd.name, cmd.desc)
	}

	for _, cmd := range commands {
		// skip commands not in this group
		if !strings.Contains(cmd.name, group) || cmd.name == group {
			continue
		}

		rest := strings.Replace(cmd.name, group+" ", "", 1)
		if strings.Contains(rest, " ") {
			continue
		}

		pErr("    %-10.10s%s\n", rest, cmd.desc)
	}
}

func IdentifyCommand(args []string) (string, []string) {
	cArgs := args
	for len(cArgs) > 0 {
		name := strings.Join(cArgs, " ")
		cArgs = cArgs[:len(cArgs)-1]

		if _, ok := commands[name]; ok {
			return name, args[len(cArgs)+1:]
		}
	}

	return "data", args[len(cArgs)+1:]
}

func DispatchCommand(name string, args []string) {
	dErr("dispatching command %s\n", name)

	cmd, ok := commands[name]
	if ok {
		err := cmd.handler(args)
		if err != nil {
			pErr("data %s: %s\n", name, err)
		}
	} else {
		pErr("data: unknown command \"%s\"\n", name)
		pErr("Run `data help` for usage.\n")
	}
}
