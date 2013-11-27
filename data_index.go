package data

import (
	"fmt"
)

const MainDataIndexURL = "http://datadex.io"

var MainDataIndex = &DataIndex{URL: MainDataIndexURL}

type DataIndex struct {
	URL string
}

func (i *DataIndex) ArchiveURL(h *Handle) string {
	ref := "master"
	if len(h.Version) > 0 {
		ref = h.Version
	}
	return fmt.Sprintf("%s/%s/archive/%s%s", i.URL, h.Path, ref, ArchiveSuffix)
}
