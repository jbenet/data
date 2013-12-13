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
	RC("get", "Download and install dataset.", getCmd)
	RC("list", "List installed datasets.", listCmd)
	RC("info", "Show dataset information.", infoCmd)
	RC("help", "Show usage information.", helpCmd)
	RC("version", "Show data version information.", versionCmd)
	RC("upload", "Upload dataset to storage service.", uploadCmd)
	RC("manifest", "Generate dataset manifest.", manifestCmd)
}

func PrintCommands() {
	for _, cmd := range commands {
		pErr("    %-10.10s%s\n", cmd.name, cmd.desc)
	}
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
