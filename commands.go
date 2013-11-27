package data

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
	RC("get", "Download and install dataset.", GetCmd)
	RC("list", "List installed datasets.", ListCmd)
	RC("info", "Show dataset information.", InfoCmd)
	RC("help", "Show usage information.", HelpCmd)
	RC("version", "Show data version information.", VersionCmd)
}

func PrintCommands() {
	for _, cmd := range commands {
		Err("    %-10.10s%s\n", cmd.name, cmd.desc)
	}
}

func DispatchCommand(name string, args []string) {
	DErr("dispatching command %s\n", name)

	cmd, ok := commands[name]
	if ok {
		err := cmd.handler(args)
		if err != nil {
			Err("data %s: %s\n", name, err)
		}
	} else {
		Err("data: unknown command \"%s\"\n", name)
		Err("Run `data help` for usage.\n")
	}
}
