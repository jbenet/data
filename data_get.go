package data

import (
	"fmt"
	"github.com/jbenet/commander"
	"os"
	"path"
	"strings"
)

var cmd_data_get = &commander.Command{
	UsageLine: "get [<dataset>|<url>]",
	Short:     "Download and install dataset.",
	Long: `data get - Download and install dataset.

    Downloads the dataset specified, and installs its files into the
    current dataset working directory.

    The dataset argument can be any of:

    HANDLE: Handle of the form <author>/<name>[.<fmt>][@<ref>].
            Looks up handle on the specified (default) datadex.

    URL:    Direct url to any dataset on any datadex. (TODO)

    PATH:   Filesystem path to any locally installed dataset. (TODO)


    Loosely, data-get's process is:

    - Locate dataset Datafile and Manifest. (via provided argument).
    - Download Datafile and Manifest, to local Repository.
    - Download Blobs, listed in Manifest to local Repository.
    - Reconstruct Files, listed in Manifest.
    - Install Files, into working directory.

  `,
	Run: getCmd,
}

func getCmd(c *commander.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("%v requires a <dataset> argument.", c.FullName())
	}

	return GetDataset(args[0])
}

func GetDataset(dataset string) error {
	dataset = strings.ToLower(dataset)

	// add lookup in datadex here.
	h := NewHandle(dataset)
	if h.Valid() {
		return GetDatasetFromIndex(h)
	}

	return fmt.Errorf("Unclear how to handle dataset identifier: %s", dataset)
}

func GetDatasetFromIndex(h *Handle) error {
	di, err := NewMainDataIndex()
	if err != nil {
		return err
	}

	pOut("Downloading %s from %s.\n", h.Dataset(), mainIndexName)

	// Prepare local directories
	dir := path.Join(DatasetDir, h.Path())
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	if err := os.Chdir(dir); err != nil {
		return err
	}

	// download manifest
	if err := di.downloadManifest(h); err != nil {
		return err
	}

	// download pack
	p, err := NewPack()
	if err != nil {
		return err
	}

	if err := p.Download(); err != nil {
		return err
	}

	df, err := NewDatafile(DatafileName)
	if err != nil {
		return err
	}

	pOut("\nInstalled %s at %s\n", df.Dataset, dir)
	return nil
}

func (d *DataIndex) downloadManifest(h *Handle) error {
	v := h.Version
	if len(v) == 0 {
		v = RefLatest
	}

	ri := d.RefIndex(h.Path())
	ref, err := ri.VersionRef(v)
	if err != nil {
		return fmt.Errorf("Error finding version %s. %s", v, err)
	}

	err = d.getBlob(ref, ManifestFileName)
	if err != nil {
		return err
	}

	return nil
}
