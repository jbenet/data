package data

import (
	"fmt"
	"github.com/jbenet/commander"
	"strings"
)

var cmd_data_upload = &commander.Command{
	UsageLine: "upload",
	Short:     "Upload dataset to storage service.",
	Long: `data upload - Upload dataset to storage service.

    Uploads dataset to a storage service. This means uploading the
    dataset's data-created repository (blobs) to a remote blobstore,
    part of the storage service. Making all the blobs representing
    the dataset (including the Datafile and Manifest) accessible
    through the storage service.

    Uploading IS NOT publishing. Uploading is a (required) step of
    publishing (see data-publish). Usually, users don't need to run
    data-upload; they instead mean to use data-publish.

    WARNING: most of what this command does will be subsumed by:
    - data-blob (to get/put blobs)
    - data-package (to fill out Datafile, build manifest, etc)
  `,
	Run: uploadCmd,
}

func uploadCmd(c *commander.Command, args []string) error {
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
