package data

import (
	"fmt"
)

type DataIndex struct {
	Url       string
	BlobStore blobStore
}

var mainDataIndex *DataIndex

func (i *DataIndex) ArchiveUrl(h *Handle) string {
	ref := "master"
	if len(h.Version) > 0 {
		ref = h.Version
	}
	return fmt.Sprintf("%s/%s/archive/%s%s", i.Url, h.Path(), ref, ArchiveSuffix)
}

// why not use `func init()`? some commands don't need an index
// is annoying to error out on an S3 key when S3 isn't needed.
func NewMainDataIndex() (*DataIndex, error) {
	if mainDataIndex != nil {
		return mainDataIndex, nil
	}

	blobStore, err := NewS3Store("datadex.archives")
	if err != nil {
		return nil, err
	}

	mainDataIndex := &DataIndex{Url: "http://datadex.io"}
	mainDataIndex.BlobStore = blobStore
	return mainDataIndex, nil
}
