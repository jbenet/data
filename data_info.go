package main

import (
	"fmt"
)

func InfoCmd(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("info requires an argument.")
	}

	return DatasetInfo(args[0])
}

func DatasetInfo(dataset string) error {
	df, err := NewDatafile(dataset)
	if err != nil {
		return fmt.Errorf("Invalid dataset handle: %s", dataset)
	}

	buf, err := df.Marshal()
	DOut("%s\n", buf)

	return nil
}
