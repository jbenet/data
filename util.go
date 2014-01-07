package data

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"os/exec"
	"unicode"
)

var Debug bool

// Shorthand printing functions.
func pErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func pOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func dErr(format string, a ...interface{}) {
	if Debug {
		pErr(format, a...)
	}
}

func dOut(format string, a ...interface{}) {
	if Debug {
		pOut(format, a...)
	}
}

// Checks whether string is a hash (sha1)
func isHash(hash string) bool {
	if len(hash) != 40 {
		return false
	}

	for _, r := range hash {
		if !unicode.Is(unicode.ASCII_Hex_Digit, r) {
			return false
		}
	}

	return true
}

func shortHash(hash string) string {
	return hash[:7]
}

func readerHash(r io.Reader) (string, error) {
	bf := bufio.NewReader(r)
	h := sha1.New()
	_, err := bf.WriteTo(h)
	if err != nil {
		return "", err
	}

	hex := fmt.Sprintf("%x", h.Sum(nil))
	return hex, nil
}

func copyFile(src string, dst string) error {
	cmd := exec.Command("cp", src, dst)
	return cmd.Run()
}

func set(slice []string) []string {
	dedup := []string{}
	elems := map[string]bool{}
	for _, elem := range slice {
		_, seen := elems[elem]
		if !seen {
			dedup = append(dedup, elem)
			elems[elem] = true
		}
	}
	return dedup
}

func validHashes(hashes []string) (valid []string, err error) {
	hashes = set(hashes)

	// append only valid hashes
	for _, hash := range hashes {
		if isHash(hash) {
			valid = append(valid, hash)
		} else {
			err = fmt.Errorf("invalid <hash>: %v", hash)
		}
	}

	return
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
