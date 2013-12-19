package data

import (
	"github.com/jbenet/commander"
)

var Cmd_data = &commander.Command{
	UsageLine: "data [<flags>] <command> [<args>]",
	Short:     "dataset package manager",
	Long: `data - dataset package manager
  `,
	Subcommands: []*commander.Command{
		cmd_data_version,
		cmd_data_info,
		cmd_data_list,
		cmd_data_get,
		cmd_data_manifest,
		cmd_data_pack,
		cmd_data_blob,
	},
}
