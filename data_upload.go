package data

import (
	"fmt"
	"strings"
)

func uploadCmd(args []string) error {
	if len(args) < 1 {
		return UploadDataset("datadex")
	} else {
		return UploadDataset(args[0])
	}
}

func UploadDataset(service string) error {

	// ensure the dataset has required information
	err := fillOutDatafileInPath(DatasetFile)
	if err != nil {
		return err
	}

	dataIndex, err := dataIndexNamed(service)
	if err != nil {
		return err
	}

	return dataIndex.uploadDataset()
}

func dataIndexNamed(name string) (*DataIndex, error) {
	switch strings.ToLower(name) {
	case "datadex":
		return mainDataIndex()
	}

	return nil, fmt.Errorf("Unsupported storage service %s", name)
}
