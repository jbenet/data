package data

import (
	"github.com/jbenet/commander"
)

const Version = "0.1.0"

var cmd_data_version = &commander.Command{
	UsageLine: "version",
	Short:     "Show data version information.",
	Long: `data version - Show data version information.

    Returns the current version of data and exits.
  `,
	Run: versionCmd,
}

func versionCmd(*commander.Command, []string) error {
	pOut("data version %s\n", Version)
	return nil
}
