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
      pack manifest   Show current package manifest.
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

  data pack manifest

    Shows the current package manifest. This may be out of date with the
    current directory contents.

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
  `,

	Subcommands: []*commander.Command{
		cmd_data_pack_make,
		cmd_data_pack_manifest,
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

var cmd_data_pack_manifest = &commander.Command{
	UsageLine: "manifest",
	Short:     "Show current package manifest.",
	Long: `data pack manifest - Show current package manifest.

    Shows the package's manifest file and exits.
    If no manifest file exists, exit with an error.

    See 'data pack'.
  `,
	Run: packManifestCmd,
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
	p, err := NewPack()
	if err != nil {
		return err
	}

	return p.GenerateFiles()
}

func packManifestCmd(c *commander.Command, args []string) error {
	p, err := NewPack()
	if err != nil {
		return err
	}

	buf, err := p.manifest.Marshal()
	if err != nil {
		return err
	}

	pOut("%s", buf)
	return nil
}

func packUploadCmd(c *commander.Command, args []string) error {
	p, err := NewPack()
	if err != nil {
		return err
	}

	if !p.manifest.Complete() {
		return fmt.Errorf(`Manifest incomplete. Before uploading, either:
      - Generate new package manifest with 'data pack make' (uses all files).
      - Finish manifest with 'data manifest' (add and hash specific files).
    `)
	}

	hashes, err := p.BlobHashes()
	if err != nil {
		return err
	}

	return putBlobs(hashes)
}

func packDownloadCmd(c *commander.Command, args []string) error {
	p, err := NewPack()
	if err != nil {
		return err
	}
	if !p.manifest.Complete() {
		return fmt.Errorf(`Manifest incomplete. Get new manifest copy.`)
	}

	hashes, err := p.BlobHashes()
	if err != nil {
		return err
	}

	return getBlobs(hashes)
}

func packCheckCmd(c *commander.Command, args []string) error {
	p, err := NewPack()
	if err != nil {
		return err
	}

	if !p.manifest.Complete() {
		pOut("Warning: manifest incomplete. Checksums may be incorrect.")
	}

	failures := 0

	for _, file := range p.manifest.AllPaths() {
		pass, err := p.manifest.Check(file)
		if err != nil {
			return err
		}

		if !pass {
			failures++
		}
	}

	count := len(*p.manifest.Files)
	if failures > 0 {
		return fmt.Errorf("data pack: %v/%v checksums failed!", failures, count)
	}

	pOut("data pack: %v checksums pass\n", count)
	return nil
}

type Pack struct {
	manifest *Manifest
	datafile *Datafile
	index    *DataIndex
}

func NewPack() (p *Pack, err error) {
	p = &Pack{}
	p.manifest = NewManifest("")

	p.datafile, _ = NewDatafile("")
	// ignore error loading datafile

	p.index, err = NewMainDataIndex()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Pack) BlobHashes() ([]string, error) {
	mfh, err := p.manifest.ManifestHash()
	if err != nil {
		return []string{}, err
	}

	hashes := p.manifest.AllHashes()
	hashes = append(hashes, mfh)
	return hashes, nil
}

func (p *Pack) GenerateFiles() error {

	// ensure the dataset has required information
	err := fillOutDatafileInPath(p.datafile.Path)
	if err != nil {
		return err
	}

	// generate manifest
	err = p.manifest.Generate()
	if err != nil {
		return err
	}

	return nil
}

// Check the blobstore to check which blobs in pack have not been uploaded.
func (p *Pack) blobsToUpload() ([]string, error) {
	missing := []string{}

	hashes, err := p.BlobHashes()
	if err != nil {
		return []string{}, err
	}

	for _, hash := range hashes {
		exists, err := p.index.hasBlob(hash)
		if err != nil {
			return []string{}, err
		}

		if !exists {
			missing = append(missing, hash)
		}
	}
	return missing, nil
}
