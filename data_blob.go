package data

import (
	"bufio"
	"fmt"
	"github.com/gonuts/flag"
	"github.com/jbenet/commander"
	"io"
	"os"
	"path"
)

var cmd_data_blob = &commander.Command{
	UsageLine: "blob <command> <hash>",
	Short:     "Manage blobs in the blobstore.",
	Long: `data blob - Manage blobs in the blobstore.

  Commands:

    put <hash> <path>     Upload blob named by <hash> to blobstore.
    get <hash> <path>     Download blob named by <hash> from blobstore.
    check <hash> <path>   Verify blob matches <hash>.
    url <hash>            Output Url for blob named by <hash>.
    show <hash>           Output blob contents for hash.
    hash <path>           Output hash for blob contents.

  Arguments:

    The <hash> argument is the blob's checksum, and id.
    The <path> argument is the blob's target file.
    If <path> is omitted, stdin/stdout are used.


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
		cmd_data_blob_url,
		cmd_data_blob_show,
		cmd_data_blob_hash,
	},
}

var cmd_data_blob_put = &commander.Command{
	UsageLine: "put <hash> <path>",
	Short:     "Upload blobs to a remote blobstore.",
	Long: `data blob put - Upload blobs to a remote blobstore.

    Upload the blob contents named by <hash> to a remote blobstore.
    Blob contents are stored locally, to be used to reconstruct files.
    In the future, the blobstore will be able to be changed. For now,
    the default blobstore/datadex is used.

    See data blob.

Arguments:

    <hash>   name (cryptographic hash, checksum) of the blob.
    <path>   path of the blob contents to upload.

  `,
	Run:  blobPutCmd,
	Flag: *flag.NewFlagSet("data-blob-put", flag.ExitOnError),
}

var cmd_data_blob_get = &commander.Command{
	UsageLine: "get <hash> [<path>]",
	Short:     "Download blobs from a remote blobstore.",
	Long: `data blob get - Download blobs from a remote blobstore.

    Download the blob contents named by <hash> from a remote blobstore.
    Blob contents are stored locally, to be used to reconstruct files.
    In the future, the blobstore will be able to be changed. For now,
    the default blobstore/datadex is used.

    See data blob.

Arguments:

    <hash>   name (cryptographic hash, checksum) of the blob.
    <path>   path to put the blob contents in.

  `,
	Run:  blobGetCmd,
	Flag: *flag.NewFlagSet("data-blob-get", flag.ExitOnError),
}

var cmd_data_blob_url = &commander.Command{
	UsageLine: "url <hash>",
	Short:     "Output Url for blob named by <hash>.",
	Long: `data blob url - Output Url for blob named by <hash>.

    Output the remote storage url for the blob contents named by <hash>.
    In the future, the blobstore will be able to be changed. For now,
    the default blobstore/datadex is used.

    See data blob.

Arguments:

    <hash>   name (cryptographic hash, checksum) of the blob.

  `,
	Run:  blobUrlCmd,
	Flag: *flag.NewFlagSet("data-blob-url", flag.ExitOnError),
}

var cmd_data_blob_show = &commander.Command{
	UsageLine: "show <hash>",
	Short:     "Output blob contents for hash.",
	Long: `data blob show - Output blob contents for hash.

    Output the blob contents stored in the blobstore for hash.
    If the blob is available locally, that copy is used (after
    hashing to verify correctness). Otherwise, it is downloaded
    from the blobstore.

    See data blob.

Arguments:

    <hash>   name (cryptographic hash, checksum) of the blob.

  `,
	Run: blobShowCmd,
}

var cmd_data_blob_hash = &commander.Command{
	UsageLine: "hash <file>",
	Short:     "Output hash for blob contents.",
	Long: `data blob hash - Output hash for blob contents.

    Output the hash of the blob contents stored in <file>

    See data blob.

Arguments:

    <file>   path of the blob contents

  `,
	Run: blobHashCmd,
}

func init() {
	cmd_data_blob.Flag.Bool("all", false, "all available blobs")
	cmd_data_blob_get.Flag.Bool("all", false, "get all available blobs")
	cmd_data_blob_put.Flag.Bool("all", false, "put all available blobs")
	cmd_data_blob_url.Flag.Bool("all", false, "urls for all available blobs")
}

type blobStore interface {
	Has(key string) (bool, error)
	Put(key string, value io.Reader) error
	Get(key string) (io.ReadCloser, error)
	Url(key string) string
}

// map { path : hash } (backward because of dup hashes)
type blobPaths map[string]string

// Handles arguments and dispatches subcommand.
func blobCmd(c *commander.Command, args []string) (blobPaths, error) {

	blobs := blobPaths{}

	// Use all blobs in the manifest if --all is passed in.
	all := c.Flag.Lookup("all").Value.Get().(bool)
	if all {
		mf := NewDefaultManifest()
		blobs = validBlobHashes(mf.Files)
		if len(blobs) < 1 {
			return nil, fmt.Errorf("%v: no blobs tracked in manifest.", c.FullName())
		}
	} else {
		switch len(args) {
		case 2:
			blobs[args[1]] = args[0]
		case 1:
			blobs[""] = args[0]
		case 0:
			return nil,
				fmt.Errorf("%v: requires <hash> argument (or --all)", c.FullName())
		}
	}

	return blobs, nil
}

func blobGetCmd(c *commander.Command, args []string) error {
	blobs, err := blobCmd(c, args)
	if err != nil {
		return err
	}
	return getBlobs(blobs)
}

func blobPutCmd(c *commander.Command, args []string) error {
	blobs, err := blobCmd(c, args)
	if err != nil {
		return err
	}
	return putBlobs(blobs)
}

func blobUrlCmd(c *commander.Command, args []string) error {
	blobs, err := blobCmd(c, args)
	if err != nil {
		return err
	}
	return urlBlobs(blobs)
}

func blobShowCmd(c *commander.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("%v: requires <hash> argument", c.FullName())
	}

	hash := args[0]
	if !IsHash(hash) {
		return fmt.Errorf("%v: invalid hash '%s'", c.FullName(), hash)
	}

	dataIndex, err := NewMainDataIndex()
	if err != nil {
		return err
	}

	return dataIndex.copyBlob(hash, os.Stdout)
}

func blobHashCmd(c *commander.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("%v: requires <file> argument", c.FullName())
	}

	hash, err := hashFile(args[0])
	if err != nil {
		return err
	}
	pOut("%s\n", hash)
	return nil
}

// Uploads all blobs to blobstore
func putBlobs(blobs blobPaths) error {
	blobs = validBlobHashes(blobs)

	dataIndex, err := NewMainDataIndex()
	if err != nil {
		return err
	}

	// flip map, to skip dupes
	flipped := map[string]string{}
	for path, hash := range blobs {
		flipped[hash] = path
	}

	for hash, path := range flipped {
		err = dataIndex.putBlob(hash, path)
		if err != nil {
			return err
		}
	}

	return nil
}

// Downloads all blobs from blobstore
func getBlobs(blobs blobPaths) error {
	blobs = validBlobHashes(blobs)

	dataIndex, err := NewMainDataIndex()
	if err != nil {
		return err
	}

	// group map, to copy dupes
	grouped := map[string][]string{}
	for path, hash := range blobs {
		g, _ := grouped[hash]
		grouped[hash] = append(g, path)
	}

	for hash, paths := range grouped {

		// download one blob
		err = dataIndex.getBlob(hash, paths[0])
		if err != nil {
			return err
		}

		// copy what we got to others
		for _, path := range paths[1:] {
			pOut("copy blob %.7s %s\n", hash, path)
			err := copyFile(paths[0], path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Shows all urls for blobs
func urlBlobs(blobs blobPaths) error {
	blobs = validBlobHashes(blobs)

	dataIndex, err := NewMainDataIndex()
	if err != nil {
		return err
	}

	for _, hash := range blobs {
		pOut("%v\n", dataIndex.urlBlob(hash))
	}

	return nil
}

// DataIndex extension to handle putting blob
func (i *DataIndex) putBlob(hash string, fpath string) error {

	// disallow empty paths
	// (stdin doesn't make sense when hashing must have already ocurred)
	if len(fpath) == 0 {
		return fmt.Errorf("put blob %.7s - error: no path supplied", hash)
	}

	fpath = path.Clean(fpath)

	// first, check the blobstore doesn't already have it.
	exists, err := i.hasBlob(hash)
	if err != nil {
		return err
	}

	if exists {
		pOut("put blob %.7s %s - exists\n", hash, fpath)
		return nil
	}

	// must verify hash before uploading (for integrity).
	// (note that there is a TOCTTOU bug here, so not safe. just helps.)
	vh, err := hashFile(fpath)
	if err != nil {
		return err
	}

	if vh != hash {
		m := "put blob: %s hash error (expected %s, got %s)"
		return fmt.Errorf(m, fpath, hash, vh)
	}

	pOut("put blob %.7s %s - uploading\n", hash, fpath)

	f, err := os.Open(fpath)
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
func (i *DataIndex) getBlob(hash string, fpath string) error {

	// disallow empty paths
	if len(fpath) == 0 {
		return fmt.Errorf("get blob %.7s - error: no path supplied", hash)
	}

	fpath = path.Clean(fpath)

	pOut("get blob %.7s %s\n", hash, fpath)
	w, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer w.Close()

	return i.copyBlob(hash, w)
}

func (i *DataIndex) copyBlob(hash string, w io.WriteCloser) error {
	r, err := i.findBlob(hash)
	if err != nil {
		return err
	}

	br := bufio.NewReader(r)
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

func (i *DataIndex) findBlob(hash string) (io.ReadCloser, error) {

	mf := NewDefaultManifest()
	paths := mf.PathsForHash(hash)
	for _, p := range paths {
		dOut("found local blob copy. verifying hash. %s\n", p)
		h, err := hashFile(p)
		if err != nil {
			continue
		}

		if hash == h {
			f, err := os.Open(p)
			if err != nil {
				continue
			}

			return f, nil
		}
	}

	dOut("no local blob copy. fetch from remote blobstore.\n")
	return i.BlobStore.Get(BlobKey(hash))
}

// DataIndex extension to check if blob exists
func (i *DataIndex) hasBlob(hash string) (bool, error) {
	return i.BlobStore.Has(BlobKey(hash))
}

// DataIndex extension to handle getting blob url
func (i *DataIndex) urlBlob(hash string) string {
	return i.BlobStore.Url(BlobKey(hash))
}

// Returns all paths associated with blob
func allBlobPaths(hash string) ([]string, error) {
	mf := NewDefaultManifest()

	paths := mf.PathsForHash(hash)

	mfh, err := mf.ManifestHash()
	if err != nil {
		return []string{}, err
	}

	if mfh == hash {
		paths = append(paths, mf.Path)
	}

	return paths, nil
}

// Returns the blobstore key for blob
func BlobKey(hash string) string {
	return fmt.Sprintf("/blob/%s", hash)
}

// Prune out invalid blob paths (bad hashes, bad paths)
func validBlobHashes(blobs blobPaths) blobPaths {
	pruned := blobPaths{}
	for fpath, hash := range blobs {
		if IsHash(hash) {
			pruned[fpath] = hash
		}
	}
	return pruned
}
