package data

import (
	"github.com/gonuts/flag"
	"github.com/jbenet/commander"
	"strings"
	"time"
)

var Cmd_data = &commander.Command{
	UsageLine: "data [<flags>] <command> [<args>]",
	Short:     "dataset package manager",
	Long: `data - dataset package manager

Basic commands:

    get         Download and install dataset.
    list        List installed datasets.
    info        Show dataset information.
    publish     Guided dataset publishing.

Tool commands:

    version     Show data version information.
    config      Manage data configuration.
    user        Manage users and credentials.
    commands    List all available commands.

Advanced Commands:

    blob        Manage blobs in the blobstore.
    manifest    Generate and manipulate dataset manifest.
    pack        Dataset packaging, upload, and download.

Use "data help <command>" for more information about a command.
`,
	Run: dataCmd,
	Subcommands: []*commander.Command{
		cmd_data_version,
		cmd_data_config,
		cmd_data_info,
		cmd_data_list,
		cmd_data_get,
		cmd_data_manifest,
		cmd_data_pack,
		cmd_data_blob,
		cmd_data_publish,
		cmd_data_user,
		cmd_data_commands,
	},
	Flag: *flag.NewFlagSet("data", flag.ExitOnError),
}

func dataCmd(c *commander.Command, args []string) error {
	pOut(c.Long)
	return nil
}

var cmd_root *commander.Command

func init() {
	// this funky alias is to resolve cyclical decl references.
	cmd_root = Cmd_data
}

var cmd_data_commands = &commander.Command{
	UsageLine: "commands",
	Short:     "List all available commands.",
	Long: `data commands - List all available commands.

    Lists all available commands (and sub-commands) and exits.
  `,
	Run: commandsCmd,
	Subcommands: []*commander.Command{
		cmd_data_commands_help,
	},
}

var cmd_data_commands_help = &commander.Command{
	UsageLine: "help",
	Short:     "List all available commands' help pages.",
	Long: `data commands help - List all available commands's help pages.

    Shows the pages of all available commands (and sub-commands) and exits.
    Outputs a markdown document, also viewable at http://datadex.io/doc/ref
  `,
	Run: commandsHelpCmd,
}

func commandsCmd(c *commander.Command, args []string) error {
	var listCmds func(c *commander.Command)
	listCmds = func(c *commander.Command) {
		pOut("%s\n", c.FullSpacedName())
		for _, sc := range c.Subcommands {
			listCmds(sc)
		}
	}

	listCmds(c.Parent)
	return nil
}

func commandsHelpCmd(c *commander.Command, args []string) error {
	pOut(referenceHeaderMsg)
	pOut("Generated on %s.\n\n", time.Now().UTC().Format("2006-01-02"))

	var printCmds func(*commander.Command, int)
	printCmds = func(c *commander.Command, level int) {
		pOut("%s ", strings.Repeat("#", level))
		pOut("%s\n\n", c.FullSpacedName())
		pOut("```\n")
		pOut("%s\n", c.Long)
		pOut("```\n\n")

		for _, sc := range c.Subcommands {
			printCmds(sc, level+1)
		}
	}

	printCmds(c.Parent.Parent, 1)
	return nil
}

const referenceHeaderMsg = `
# data command reference

This document lists every data command (including subcommands), along with
its help page. It can be viewed by running 'data commands help', and
at http://datadex.io/doc/ref

`
