package data

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

func getCmd(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("get requires a <dataset> argument.")
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
		return downloadDatasetArchive(mainDataIndex.ArchiveUrl(h))
	}

	return fmt.Errorf("Unclear how to handle dataset identifier: %s", dataset)
}

func downloadDatasetArchive(archiveUrl string) error {
	base := path.Base(archiveUrl)
	arch := path.Join(DatasetDir, ".downloads", base)

	// download the archive
	// TODO: add local caching of downloads
	pOut("Downloading archive at %s\n", archiveUrl)
	err := downloadUrlToFile(archiveUrl, arch)
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

// Url utils

const ArchiveSuffix = ".tar.gz"

func IsArchiveUrl(str string) bool {
	return isUrl(str) && strings.HasSuffix(str, ArchiveSuffix)
}

func isUrl(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func downloadUrl(url string) (*http.Response, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Got HTTP status code >= 400: %s", resp.Status)
	}

	return resp, nil
}

func urlContents(url string) ([]byte, error) {
	resp, err := downloadUrl(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func downloadUrlToFile(url string, filename string) error {
	resp, err := downloadUrl(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := createFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func createFile(filename string) (*os.File, error) {
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return file, err
}

// Extraction
func extractArchive(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	dst := strings.TrimSuffix(filename, ArchiveSuffix)
	err = os.MkdirAll(dst, 0777)
	if err != nil {
		return err
	}

	dst = path.Base(dst)
	src := path.Base(filename)
	cmd := exec.Command("tar", "xzf", src, "--strip-components", "1", "-C", dst)
	cmd.Dir = path.Dir(filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		outs := string(out)
		if strings.Contains(outs, "Error opening archive:") {
			return fmt.Errorf(outs)
		}

		return err
	}

	return nil
}
