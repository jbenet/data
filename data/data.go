package main

import (
	"os"
	"fmt"
	"github.com/jbenet/data"
	"github.com/jteeuwen/go-pkg-optarg"
)

func main() {

	optarg.UsageInfo = fmt.Sprintf("Options usage: %s [options]:", os.Args[0])
	optarg.Add("h", "help", "Show usage", false)
	optarg.Add("d", "debug", "Enter debug mode", false)
	optarg.Add("", "version", "Show version", false)
	optarg.Parse()

	forceCommand := ""
	for opt := range optarg.Parse() {
		switch opt.Name {
		case "debug":
			data.Debug = true
		case "version":
			forceCommand = "data version"
		case "help":
			forceCommand = "data help"
		}
	}

	if data.Debug {
		fmt.Fprintf(os.Stdout, "debugging on\n")
	}

	data.RegisterCommands()

	if len(forceCommand) > 0 {
		data.DispatchCommand(forceCommand, []string{})
		return
	}

	args := optarg.Remainder[:len(optarg.Remainder)/2]
	if len(args) < 1 {
		data.Usage("data")
		return
	}

	args = append([]string{"data"}, args...)
	cmd, args := data.IdentifyCommand(args)
	data.DispatchCommand(cmd, args)
}
