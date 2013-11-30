package data

import (
	"fmt"
)

func infoCmd(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("info requires an argument.")
	}

	return datasetInfo(args[0])
}

func datasetInfo(dataset string) error {
	df, err := NewDatafile(DatafilePath(dataset))
	if err != nil {
		return fmt.Errorf("Invalid dataset handle: %s", dataset)
	}

	buf, err := df.Marshal()
	if err != nil {
		return err
	}

	pOut("%s\n", buf)
	return nil
}
