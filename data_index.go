package data

import (
	"fmt"
)

type DataIndex struct {
	Url       string
	BlobStore blobStore
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
