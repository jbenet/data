package data

import (
	"fmt"
)

const mainDataIndexUrl = "http://datadex.io"

var mainDataIndex = &DataIndex{Url: mainDataIndexUrl}

type DataIndex struct {
	Url string
}

func (i *DataIndex) ArchiveUrl(h *Handle) string {
	ref := "master"
	if len(h.Version) > 0 {
		ref = h.Version
	}
	return fmt.Sprintf("%s/%s/archive/%s%s", i.Url, h.Path(), ref, ArchiveSuffix)
}
