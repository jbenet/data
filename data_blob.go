package data

import (
	"bufio"
	"fmt"
	"github.com/gonuts/flag"
	"github.com/jbenet/commander"
	"io"
	"os"
)

var cmd_data_blob = &commander.Command{
	UsageLine: "blob <command> <hash>",
	Short:     "Manage blobs in the blobstore.",
	Long: `data blob - Manage blobs in the blobstore.

    Managing blobs means:

      put <hash>    Upload blob named by <hash> to blobstore.
      get <hash>    Download blob named by <hash> from blobstore.
      check <hash>  Verify blob contents named by <hash> match <hash>.
      show <hash>   Output blob contents named by <hash>.


    What is a blob?

    Datasets are made up of files, which are made up of blobs.
    (For now, 1 file is 1 blob. Chunking to be implemented)
    Blobs are basically blocks of data, which are checksummed
    (for integrity, de-duplication, and addressing) using a crypto-
    graphic hash function (sha1, for now). If git comes to mind,
    that's exactly right.

    Local Blobstores

    data stores blobs in blobstores. Every local dataset has a
    blobstore (local caching with links TBI). Like in git, the blobs
    are stored safely in the blobstore (different directory) and can
    be used to reconstruct any corrupted/deleted/modified dataset files.

    Remote Blobstores

    data uses remote blobstores to distribute datasets across users.
    The datadex service includes a blobstore (currently an S3 bucket).
    By default, the global datadex blobstore is where things are
    uploaded to and retrieved from.

    Since blobs are uniquely identified by their hash, maintaining one
    global blobstore helps reduce data redundancy. However, users can
    run their own datadex service. (The index and blobstore are tied
    together to ensure consistency. Please do not publish datasets to
    an index if blobs aren't in that index)

    data can use any remote blobstore you wish. (For now, you have to
    recompile, but in the future, you will be able to) Just change the
    datadex configuration variable. Or pass in "-s <url>" per command.

    (data-blob is part of the plumbing, lower level tools.
    Use it directly if you know what you're doing.)
  `,

	Flag: *flag.NewFlagSet("data-blob", flag.ExitOnError),

	Subcommands: []*commander.Command{
		cmd_data_blob_put,
		cmd_data_blob_get,
	},
}

var cmd_data_blob_put = &commander.Command{
	UsageLine: "put <hash>",
	Short:     "Upload blobs to a remote blobstore.",
	Long: `data blob put - Upload blobs to a remote blobstore.

    Upload the blob contents named by <hash> to a remote blobstore.
    Blob contents are stored locally, to be used to reconstruct files.
    In the future, the blobstore will be able to be changed. For now,
    the default blobstore/datadex is used.

    See data blob.

Arguments:

    <hash>   name (cryptographic hash, checksum) of the blob.

  `,
	Run: blobPutCmd,
}

var cmd_data_blob_get = &commander.Command{
	UsageLine: "get <hash>",
	Short:     "Download blobs from a remote blobstore.",
	Long: `data blob get - Download blobs from a remote blobstore.

    Download the blob contents named by <hash> from a remote blobstore.
    Blob contents are stored locally, to be used to reconstruct files.
    In the future, the blobstore will be able to be changed. For now,
    the default blobstore/datadex is used.

    See data blob.

Arguments:

    <hash>   name (cryptographic hash, checksum) of the blob.

  `,
	Run:  blobGetCmd,
	Flag: *flag.NewFlagSet("data-blob-get", flag.ExitOnError),
}

func init() {
	cmd_data_blob.Flag.Bool("all", false, "all available blobs")
	cmd_data_blob_get.Flag.Bool("all", false, "get all available blobs")
	cmd_data_blob_put.Flag.Bool("all", false, "put all available blobs")
}

type blobStore interface {
	Put(key string, value io.Reader) error
	Get(key string) (io.ReadCloser, error)
}

// Handles arguments and dispatches subcommand.
func blobCmd(c *commander.Command, args []string) ([]string, error) {

	hashes := args

	// Use all hashes in the manifest if --all is passed in.
	all := c.Flag.Lookup("all").Value.Get().(bool)
	if all {
		mf := NewManifest("")
		hashes = mf.AllHashes()
	}

	if len(hashes) < 1 {
		return nil, fmt.Errorf("%v: requires <hash> argument (or --all)", c.FullName())
	}

	return hashes, nil
}

func blobGetCmd(c *commander.Command, args []string) error {
	hashes, err := blobCmd(c, args)
	if err != nil {
		return err
	}
	return getBlobs(hashes)
}

func blobPutCmd(c *commander.Command, args []string) error {
	hashes, err := blobCmd(c, args)
	if err != nil {
		return err
	}
	return putBlobs(hashes)
}

// Uploads all blobs named by `hashes` from blobstore
func putBlobs(hashes []string) error {

	hashes, err := validHashes(hashes)
	if err != nil {
		return err
	}

	dataIndex, err := mainDataIndex()
	if err != nil {
		return err
	}

	for _, hash := range hashes {

		paths, err := blobPaths(hash)
		if err != nil {
			return err
		}

		pOut("put blob %.7s %s\n", hash, paths[0])
		err = dataIndex.putBlob(hash, paths[0])
		if err != nil {
			return err
		}
	}

	return nil
}

// Downloads all blobs named by `hashes` from blobstore
func getBlobs(hashes []string) error {

	hashes, err := validHashes(hashes)
	if err != nil {
		return err
	}

	dataIndex, err := mainDataIndex()
	if err != nil {
		return err
	}

	for _, hash := range hashes {

		paths, err := blobPaths(hash)
		if err != nil {
			return err
		}

		// download one blob
		pOut("get blob %.7s %s\n", hash, paths[0])
		err = dataIndex.getBlob(hash, paths[0])
		if err != nil {
			return err
		}

		// copy what we got to others
		for _, path := range paths[1:] {
			pOut("get blob %.7s %s\n", hash, path)
			err := copyFile(paths[0], path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DataIndex extension to handle putting blob
func (i *DataIndex) putBlob(hash string, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bf := bufio.NewReader(f)
	err = i.BlobStore.Put(BlobKey(hash), bf)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// DataIndex extension to handle getting blob
func (i *DataIndex) getBlob(hash string, path string) error {
	r, err := i.BlobStore.Get(BlobKey(hash))
	if err != nil {
		return err
	}
	defer r.Close()

	br := bufio.NewReader(r)
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, br)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = r.Close()
	if err != nil {
		return err
	}

	return nil
}

// Returns all paths associated with blob
func blobPaths(hash string) ([]string, error) {
	mf := NewManifest("")
	return mf.PathsForHash(hash)
}

// Returns the blobstore key for blob
func BlobKey(hash string) string {
	return fmt.Sprintf("/blob/%s", hash)
}
