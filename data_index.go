package data

import (
	"bufio"
	"fmt"
	"os"
)

type DataIndex struct {
	Url       string
	BlobStore kvStore
}

func (i *DataIndex) ArchiveUrl(h *Handle) string {
	ref := "master"
	if len(h.Version) > 0 {
		ref = h.Version
	}
	return fmt.Sprintf("%s/%s/archive/%s%s", i.Url, h.Path(), ref, ArchiveSuffix)
}

func mainDataIndex() (*DataIndex, error) {
	blobStore, err := NewS3Store("datadex.archives")
	if err != nil {
		return nil, err
	}

	i := &DataIndex{Url: "http://datadex.io"}
	i.BlobStore = blobStore
	return i, nil
}

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
		err = i.uploadBlob(f, h)
		if err != nil {
			return err
		}
	}

	// upload manifest itself
	mfh, err := hashFile(mf.Path)
	if err != nil {
		return err
	}

	err = i.uploadBlob(mf.Path, mfh)
	if err != nil {
		return err
	}

	return nil
}

func (i *DataIndex) uploadBlob(path string, hash string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bf := bufio.NewReader(f)
	err = i.BlobStore.Put(hash, bf)
	return err
}
