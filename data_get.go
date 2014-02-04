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
	var datasets []string

	if len(args) > 0 {
		// if args, get those datasets.
		datasets = args
	} else {
		// if no args, use Datafile dependencies
		df, _ := NewDefaultDatafile()
		for _, dep := range df.Dependencies {
			if NewHandle(dep).Valid() {
				datasets = append(datasets, dep)
			}
		}
	}

	if len(datasets) == 0 {
		return fmt.Errorf("%v: no datasets specified.\nEither enter a <dataset> "+
			"argument, or add dependencies in a Datafile.", c.FullName())
	}

	for _, ds := range datasets {
		err := GetDataset(ds)
		if err != nil {
			return err
		}
	}

	if len(datasets) == 0 {
		return nil
	}

	// If many, Installation Summary
	pErr("---------\n")
	for _, ds := range datasets {
		err := installedDatasetMessage(ds)
		if err != nil {
			pErr("%v\n", err)
		}
	}
	return nil
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

	pErr("Downloading %s from %s (%s).\n", h.Dataset(), di.Name, di.Http.Url)

	// Get manifest ref
	mref, err := di.handleRef(h)
	if err != nil {
		return err
	}

	// Prepare local directories
	dir := path.Join(DatasetDir, h.Path())
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(dir); err != nil {
		return err
	}

	// move back out
	defer os.Chdir(cwd)

	// download manifest
	if err := downloadManifest(di, mref); err != nil {
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

	pErr("\n")
	return nil
}

func (d *DataIndex) handleRef(h *Handle) (string, error) {
	v := h.Version
	if len(v) == 0 {
		v = RefLatest
	}

	ri := d.RefIndex(h.Path())
	ref, err := ri.VersionRef(v)
	if err != nil {
		if strings.Contains(err.Error(), "404 page not found") {
			return "", fmt.Errorf("Error: %v not found.", h.Dataset())
		}
		return "", fmt.Errorf("Error finding manifest for %v. %s", h.Dataset(), err)
	}

	return ref, nil
}

func downloadManifest(d *DataIndex, ref string) error {
	return d.getBlob(ref, ManifestFileName)
}

func installedDatasetMessage(dataset string) error {
	h := NewHandle(dataset)
	fpath := DatafilePath(h.Path())
	df, err := NewDatafile(fpath)
	if err != nil {
		return err
	}

	pOut("Installed %s at %s\n", df.Dataset, path.Dir(fpath))
	return nil
}
