package data

import (
	"fmt"
	"path"
	"strings"
)

// <author>/<name>[.<format>][@<tag>]

type Handle struct {
	Author  string
	Name    string
	Format  string
	Version string
}

// There are problems with goyaml setters/getters.
// Unmarshaling fails.
//
// func (d Handle) GetYAML() (string, interface{}) {
// 	pOut("GetYAML\n")
// 	return "", d.string
// }
//
// func (d Handle) SetYAML(tag string, value interface{}) bool {
// 	s, ok := value.(string)
// 	d.string = s
// 	pOut("SetYAML %s %s\n", d.string, &d)
// 	return ok
// }

func NewHandle(s string) *Handle {
	d := new(Handle)
	d.SetDataset(s)
	return d
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
func (d *Handle) SetDataset(s string) {

	nam_idx := strings.Index(s, "/")
	if nam_idx < 0 {
		nam_idx = 0
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
}

func (d *Handle) GoString() string {
	return d.Dataset()
}

func (d *Handle) Valid() bool {
	return IsDatasetHandle(d.Dataset())
}

// utils

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

func handleError(handle string, problem string) error {
	return fmt.Errorf("Invalid handle (%s): %s", problem, handle)
}

func IsDatasetHandle(str string) bool {
	return HandleRegexp.MatchString(str)
}
