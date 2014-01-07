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

    URL:    Direct url to any dataset on any datadex.

    PATH:   Filesystem path to any locally installed dataset.


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

	if IsArchiveUrl(dataset) {
		return downloadDatasetArchive(dataset)
	}

	// add lookup in datadex here.
	h := NewHandle(dataset)
	if h.Valid() {
		dataIndex, err := NewMainDataIndex()
		if err != nil {
			return err
		}

		return downloadDatasetArchive(dataIndex.ArchiveUrl(h))
	}

	return fmt.Errorf("Unclear how to handle dataset identifier: %s", dataset)
}

func downloadDatasetArchive(archiveUrl string) error {
	base := path.Base(archiveUrl)
	arch := path.Join(DatasetDir, ".downloads", base)

	// download the archive
	// TODO: add local caching of downloads
	pOut("Downloading archive at %s\n", archiveUrl)
	err := httpWriteToFile(archiveUrl, arch)
	if err != nil {
		return err
	}

	// untar the archive
	dOut("Extracting archive at %s\n", arch)
	err = extractArchive(arch)
	if err != nil {
		return err
	}

	// find place from Datafile
	arch_dir := strings.TrimSuffix(arch, ArchiveSuffix)
	df, err := NewDatafile(path.Join(arch_dir, DatasetFile))
	if err != nil {
		return err
	}
	pOut("%s downloaded\n", df.Dataset)

	// move into place
	new_path := path.Join(DatasetDir, df.Dataset)
	err = os.MkdirAll(path.Dir(new_path), 0777)
	if err != nil {
		return err
	}

	_, err = os.Stat(new_path)
	if err == nil {
		return fmt.Errorf("error: dataset already installed at %s\n"+
			"Remove and try again.\n", new_path)
	}

	err = os.Rename(arch_dir, new_path)
	if err != nil {
		return err
	}
	pOut("%s installed\n", df.Dataset)

	return nil
}
