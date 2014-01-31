package data

import (
	"github.com/gonuts/flag"
	"github.com/jbenet/commander"
)

const Version = "0.1.1"

var cmd_data_version = &commander.Command{
	UsageLine: "version",
	Short:     "Show data version information.",
	Long: `data version - Show data version information.

    Returns the current version of data and exits.
  `,
	Run:  versionCmd,
	Flag: *flag.NewFlagSet("data-user-auth", flag.ExitOnError),
}

func init() {
	cmd_data_version.Flag.Bool("number", false, "show only the number")
}

func versionCmd(c *commander.Command, _ []string) error {
	number := c.Flag.Lookup("number").Value.Get().(bool)
	if !number {
		pOut("data version ")
	}
	pOut("%s\n", Version)
	return nil
}
