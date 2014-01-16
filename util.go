package data

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"github.com/dotcloud/docker/pkg/term"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

var Debug bool
var NotImplementedError = fmt.Errorf("Error: not implemented yet.")

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
func IsHash(hash string) bool {
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

func StringHash(s string) (string, error) {
	r := strings.NewReader(s)
	h := sha1.New()
	_, err := r.WriteTo(h)
	if err != nil {
		return "", err
	}

	hex := fmt.Sprintf("%x", h.Sum(nil))
	return hex, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return readerHash(f)
}

func catFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	_, err = io.Copy(os.Stdout, br)
	return err
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
		if IsHash(hash) {
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

func httpExists(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	c := resp.StatusCode
	switch {
	case 200 <= c && c < 400:
		return true, nil
	case 400 <= c && c < 500:
		return false, nil
	default:
		return false, fmt.Errorf("Network or server error retrieving: %s", url)
	}
}

func httpGet(url string) (*http.Response, error) {
	dOut("http get %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	c := resp.StatusCode
	if 200 <= c && c < 400 {
		return resp, nil
	}

	e, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	s := strings.TrimSpace(string(e[:]))
	return nil, fmt.Errorf("HTTP error status code: %d (%s)", c, s)
}

func httpPost(url string, bt string, b io.Reader) (*http.Response, error) {
	dOut("http post %s\n", url)
	resp, err := http.Post(url, bt, b)
	if err != nil {
		return nil, err
	}

	c := resp.StatusCode
	if 200 <= c && c < 400 {
		return resp, nil
	}

	e, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	s := strings.TrimSpace(string(e[:]))
	return nil, fmt.Errorf("HTTP error status code: %d (%s)", c, s)
}

func httpReadAll(url string) ([]byte, error) {
	resp, err := httpGet(url)
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

func httpWriteToFile(url string, filename string) error {
	resp, err := httpGet(url)
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

// Return array of all files matching func
func FindFiles(dir string, recursive bool, skipHidden bool) ([]string, error) {
	filenames := []string{}
	walkFn := func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {

			// entirely skip hidden dirs
			if skipHidden &&
				len(info.Name()) > 1 && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			if !recursive {
				return filepath.SkipDir
			}

			return nil
		}

		// hidden?
		if skipHidden && strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		filenames = append(filenames, path)
		return nil
	}

	err := filepath.Walk(dir, walkFn)
	if err != nil {
		return nil, err
	}

	return filenames, nil
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

// Input
func readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(line), nil
}

func readInputSilent() (string, error) {
	fd := os.Stdin.Fd()
	s, _ := term.SaveState(fd)
	term.DisableEcho(fd, s)

	input, err := readInput()
	term.RestoreTerminal(fd, s)

	pOut("\n")
	return input, err
}

// Exec helper
func execCmdArgs(path string, args []string) (string, []string) {
	if args == nil {
		args = []string{}
	}

	parts := strings.Split(path, " ")
	if len(parts) > 1 {
		path = parts[0]
		args = append(parts[1:], args...)
	}

	return path, args
}

// Map sorting -- lifted from
// https://groups.google.com/d/msg/golang-nuts/FT7cjmcL7gw/Gj4_aEsE_IsJ

// A data structure to hold a key/value pair.
type pair struct {
	Key   string
	Value string
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type pairList []pair

func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]string) pairList {
	p := make(pairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = pair{k, v}
	}
	sort.Sort(p)
	return p
}
