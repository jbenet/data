package data

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DataManifest = ".data-manifest"
const noHash = "h"

func manifestCmd(args []string) error {
	return generateManifest()
}

func generateManifest() error {

	mf := NewManifest(DataManifest)

	// add new files to manifest file
	// (for now add everything. `data manifest {add,rm}` in future)
	for _, f := range listAllFiles(".") {
		mf.Add(f)
	}

	// warn about manifest-listed files missing from directory
	// (basically, missing things. User removes individually, or `rm --missing`)

	// Once all files are listed, hash all the files, storing the hashes.
	for f, h := range *mf.Files {
		if h != noHash {
			continue
		}

		err := mf.Hash(f)
		if err != nil {
			return err
		}
	}

	return nil
}

type Manifest struct {
	file  "-"
	Files *map[string]string ""
}

func NewManifest(path string) *Manifest {
	mf := &Manifest{file: file{Path: path}}

	// initialize map
	mf.Files = &map[string]string{}
	mf.file.format = mf.Files

	// attempt to load
	mf.ReadFile()
	return mf
}

func (mf *Manifest) Add(path string) {
	// check, dont override (could have hash value)
	_, exists := (*mf.Files)[path]
	if !exists {
		(*mf.Files)[path] = noHash
		pOut("data manifest: added %s\n", path)
	}
}

func (mf *Manifest) Hash(path string) error {
	h, err := hashFile(path)
	if err != nil {
		return err
	}

	(*mf.Files)[path] = h

	// Write out file (store incrementally)
	err = mf.WriteFile()
	if err != nil {
		return err
	}

	pOut("data manifest: hashed %.7s %s\n", h, path)
	return nil
}

func listAllFiles(path string) []string {

	files := []string{}
	walkFn := func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {

			// entirely skip hidden dirs
			if len(info.Name()) > 1 && strings.HasPrefix(info.Name(), ".") {
				dOut("data manifest: skipping %s/\n", info.Name())
				return filepath.SkipDir
			}

			// skip datasets/
			if path == DatasetDir {
				dOut("data manifest: skipping %s/\n", info.Name())
				return filepath.SkipDir
			}

			// dont store dirs
			return nil
		}

		// skip manifest file
		if path == DataManifest {
			dOut("data manifest: skipping %s\n", info.Name())
			return nil
		}

		files = append(files, path)
		return nil
	}

	filepath.Walk(path, walkFn)
	return files
}

func hashFile(path string) (string, error) {

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bf := bufio.NewReader(f)
	h := sha1.New()
	_, err = bf.WriteTo(h)
	if err != nil {
		return "", err
	}

	hex := fmt.Sprintf("%x", h.Sum(nil))
	return hex, nil
}
