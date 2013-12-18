package data

import (
	"fmt"
	"github.com/jbenet/commander"
)

var cmd_data_info = &commander.Command{
	UsageLine: "info <dataset>",
	Short:     "Show dataset information.",
	Long: `data info - Show dataset information.

    Returns the Datafile corresponding to <dataset> and exits.
  `,
	Run: infoCmd,
}

func infoCmd(c *commander.Command, args []string) error {
	if len(args) < 1 {

		return fmt.Errorf("%v requires a <dataset> argument.", c.FullName())
	}

	return datasetInfo(args[0])
}

func datasetInfo(dataset string) error {
	df, err := NewDatafile(DatafilePath(dataset))
	if err != nil {
		dErr("Error: %s\n", err)
		return fmt.Errorf("Invalid dataset handle: %s", dataset)
	}

	buf, err := df.Marshal()
	if err != nil {
		return err
	}

	pOut("%s\n", buf)
	return nil
}
