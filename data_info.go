package data

import (
	"fmt"
	"github.com/jbenet/commander"
)

var cmd_data_info = &commander.Command{
	UsageLine: "info [<dataset>]",
	Short:     "Show dataset information.",
	Long: `data info - Show dataset information.

    Returns the Datafile corresponding to <dataset> (or in current
    directory) and exits.
  `,
	Run: infoCmd,
}

func infoCmd(c *commander.Command, args []string) error {
	if len(args) < 1 {
		return datasetInfo(DatafileName)
	}

	return datasetInfo(DatafilePath(args[0]))
}

func datasetInfo(path string) error {
	df, err := NewDatafile(path)
	if err != nil {
		dErr("Error: %s\n", err)
		return fmt.Errorf("Invalid dataset path: %s", path)
	}

	buf, err := df.Marshal()
	if err != nil {
		return err
	}

	pOut("%s\n", buf)
	return nil
}
