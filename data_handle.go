package data

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// <author>/<name>[.<format>][@<tag>]
type Handle struct {
	Author  string "-"
	Name    string "-"
	Version string "-"
	Format  string "-"
}

func NewHandle(s string) (*Handle, error) {
	d := new(Handle)

	err := d.SetString(s)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Handle) Dataset() string {
	s := d.Path()

	if len(d.Format) > 0 {
		s = fmt.Sprintf("%s.%s", s, d.Format)
	}

	if len(d.Version) > 0 {
		s = fmt.Sprintf("%s@%s", s, d.Version)
	}

	return s
}

func (d *Handle) Path() string {
	return path.Join(d.Author, d.Name)
}

// order: rsplit @, split /, rsplit .
func (d *Handle) SetString(s string) error {

	nam_idx := strings.Index(s, "/")
	if nam_idx < 0 {
		return handleError(s, "no author/name separator")
	}

	ver_idx := strings.LastIndex(s, "@")
	if ver_idx < 0 {
		ver_idx = len(s) // no version in handle.
	}

	// this precludes names that have periods... use different delimiter?
	fmt_idx := strings.LastIndex(s[nam_idx:ver_idx], ".")
	if fmt_idx < 0 {
		fmt_idx = ver_idx // no format in handle.
	}

	// parts
	d.Author = slice(s, 0, nam_idx)
	d.Name = slice(s, nam_idx, fmt_idx)
	d.Format = slice(s, fmt_idx+1, ver_idx)
	d.Version = slice(s, ver_idx+1, len(s))
	return nil
}

func slice(s string, from int, to int) string {
	from = maxInt(from, 0)
	to = minInt(to, len(s))
	return s[minInt(from, to):to]
}

// https://groups.google.com/forum/#!topic/golang-nuts/dbyqx_LGUxM is silly.
func minInt(x, y int) (r int) {
	if x < y {
		return x
	}
	return y
}

func maxInt(x, y int) (r int) {
	if x > y {
		return x
	}
	return y
}

func (d *Handle) GoString() string {
	return d.Dataset()
}

func (d *Handle) GetYAML() (tag string, value interface{}) {
	pOut("GetYAML called\n")
	return "", d.Dataset()
}

func (d *Handle) SetYAML(tag string, value interface{}) bool {
	pOut("SetYAML called\n")

	str, ok := value.(string)
	if !ok {
		return false
	}

	err := d.SetString(str)
	return err == nil
}

func handleError(handle string, problem string) error {
	return fmt.Errorf("Invalid handle (%s): %s", problem, handle)
}

const identRE = "[a-z0-9_-.]+"
const namRE = identRE
const autRE = identRE
const fmtRE = identRE
const refRE = identRE
const hdlRE = "((" + namRE + ")/(" + autRE + "))(\\." + fmtRE + ")?" + "(@" + refRE + ")?"

func IsDatasetHandle(str string) bool {
	match, err := regexp.MatchString(hdlRE, str)
	return err != nil && match
}
