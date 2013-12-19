package data

import (
	"fmt"
	"github.com/jbenet/commander"
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
	return UploadDataset(args)
}

func (i *DataIndex) uploadFileOrBlob(pathOrHash string) error {
	mf := NewManifest("")
	h, p, err := mf.Pair(pathOrHash)
	if err != nil {
		return err
	}

	return i.putBlob(h, p)
}

// func (i *DataIndex) downloadFileOrBlob(pathOrHash string) error {
// 	mf := NewManifest("")
// 	h, p, err := mf.Pair(pathOrHash)
// 	if err != nil {
// 		return err
// 	}

// 	return i.getBlob(h, p)
// }

func UploadDataset(args []string) error {

	// ensure the dataset has required information
	err := fillOutDatafileInPath(DatasetFile)
	if err != nil {
		return err
	}

	dataIndex, err := mainDataIndex()
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return dataIndex.uploadDataset()
	}

	for _, arg := range args {
		err := dataIndex.uploadFileOrBlob(arg)
		if err != nil {
			return err
		}
	}

	return nil
}

// func DownloadDataset(args []string) error {
// 	dataIndex, err := mainDataIndex()
// 	if err != nil {
// 		return err
// 	}

// 	if len(args) < 1 {
// 		return dataIndex.downloadDataset()
// 	}

// 	for _, arg := range args {
// 		err := dataIndex.downloadFileOrBlob(arg)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (i *DataIndex) uploadDataset() error {
	err := i.uploadBlobs()
	if err != nil {
		return err
	}

	// upload manifest + datafile to index
	return fmt.Errorf("upload manifest + datafile to index not implemented\n")
	return nil
}

func (i *DataIndex) uploadBlobs() error {
	// regenerate manifest
	mf, err := NewGeneratedManifest("")
	if err != nil {
		return err
	}

	// upload manifest files
	for f, h := range *mf.Files {
		err = i.putBlob(h, f)
		if err != nil {
			return err
		}
	}

	// upload manifest itself
	mfh, err := hashFile(mf.Path)
	if err != nil {
		return err
	}

	err = i.putBlob(mfh, mf.Path)
	if err != nil {
		return err
	}

	return nil
}
