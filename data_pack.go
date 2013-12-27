package data

import (
	"fmt"
	"github.com/jbenet/commander"
)

var cmd_data_pack = &commander.Command{
	UsageLine: "pack [ download | upload ]",
	Short:     "Dataset packaging, upload, and download.",
	Long: `data pack - Dataset packaging, upload, and download.

    Commands:

      pack make       Create or update package description.
      pack upload     Upload package to remote storage.
      pack download   Download package from remote storage.
      pack checksum   Verify all file checksums match.


    What is a data package?

    A data package represents a single dataset, a unit of information.
    data makes it easy to find, download, create, publish, and maintain
    these datasets/packages.

    Dataset packages are simply file directories with two extra files:
    - Datafile, containing dataset description and metadata
    - Manifest, containing dataset file paths and checksums
    (See 'data help datafile' and 'data help manifest'.)

    data pack make

    'Packing' is the process of generating the package's Datafile and
    Manifest. The Manifest is built automatically, but the Datafile
    requires user input, to specify name, author, description, etc.

    data pack upload

    Packages, once 'packed' (Datafile + Manifest created), can be uploaded
    to a remote storage service (by default, the datadex). This means
    uploading all the package's files (blobs) not already present in the
    storage service. This is determined using a checksum.

    data pack download

    Similarly, packages can be downloaded or reconstructed in any directory
    from the Datafile and Manifest. Running 'data pack download' ensures
    all files listed in the Manifest are downloaded to the directory.

    data pack checksum

    Packages can be verified entirely by calling the 'data pack checksum'
    command. It re-hashes every file and ensures the checksums match.

    Packages can be published to the dataset index using 'data publish'.
  `,

	Subcommands: []*commander.Command{
		cmd_data_pack_make,
		cmd_data_pack_upload,
		cmd_data_pack_download,
		cmd_data_pack_check,
	},
}

var cmd_data_pack_make = &commander.Command{
	UsageLine: "make",
	Short:     "Create or update package description.",
	Long: `data pack upload - Upload package contents to remote storage.

    Makes the package's description files:
    - Datafile, containing dataset description and metadata (prompts)
    - Manifest, containing dataset file paths and checksums (generated)

    See 'data pack'.
  `,
	Run: packMakeCmd,
}

var cmd_data_pack_upload = &commander.Command{
	UsageLine: "upload",
	Short:     "Upload package contents to remote storage.",
	Long: `data pack upload - Upload package contents to remote storage.

    Uploads package's files (blobs) to a remote storage service (datadex).
    Blobs are named by their hash (checksum), so data can deduplicate.
    Meaning, data can easily tell whether the service already has each
    file, avoiding redundant uploads, saving bandwidth, and leveraging
    the data uploaded along with other datasets.

    See 'data pack'.
  `,
	Run: packUploadCmd,
}

var cmd_data_pack_download = &commander.Command{
	UsageLine: "download",
	Short:     "Download package contents from remote storage.",
	Long: `data pack download - Download package contents from remote storage.

    Downloads package's files (blobs) from remote storage service (datadex).
    Blobs are named by their hash (checksum), so data can deduplicate and
    ensure integrity. Meaning, data can avoid redundant downloads, saving
    bandwidth and speed, as well as verify the correctness of files with
    their checksum, preventing corruption.

    See 'data pack'.
  `,
	Run: packDownloadCmd,
}

var cmd_data_pack_check = &commander.Command{
	UsageLine: "check",
	Short:     "Verify all file checksums match.",
	Long: `data pack check - Verify all file checksums match.

    Verifies all package's file (blob) checksums match hashes stored in
    the Manifest. This is the way to check package-wide integrity. If any
    checksums FAIL, it is suggested that the files be re-downloaded (using
    'data pack download' or 'data blob get <hash>').

    See 'data pack'.
  `,
	Run: packCheckCmd,
}

func packMakeCmd(c *commander.Command, args []string) error {
	_, err := packGenerateFiles()
	return err
}

func packUploadCmd(c *commander.Command, args []string) error {
	mf := NewManifest("")
	if len(*mf.Files) < 1 {
		return fmt.Errorf("No files in manifest. " +
			"Generate manifest with 'data pack make'")
	}
	return putBlobs(mf.AllHashes())
}

func packDownloadCmd(c *commander.Command, args []string) error {
	mf := NewManifest("")
	return getBlobs(mf.AllHashes())
}

func packCheckCmd(c *commander.Command, args []string) error {
	count := 0

	mf := NewManifest("")
	for _, file := range mf.AllPaths() {
		err := mf.Check(file)
		if err != nil {
			count++
		}
	}

	if count > 0 {
		return fmt.Errorf("data pack: %v checksums failed!", count)
	}

	pOut("data pack: %v checksums pass\n", len(*mf.Files))
	return nil
}

func packGenerateFiles() (*Manifest, error) {

	// ensure the dataset has required information
	err := fillOutDatafileInPath(DatasetFile)
	if err != nil {
		return nil, err
	}

	// regenerate manifest
	mf, err := NewGeneratedManifest("")
	if err != nil {
		return nil, err
	}

	return mf, nil
}
