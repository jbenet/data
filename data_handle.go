package main

import (
	"fmt"
	"strings"
)

// <author>/<name>[.<format>][@<tag>]
type Handle struct {
	Handle string

	Author  string "-"
	Name    string "-"
	Version string "-"
	Format  string "-"

	Path string "-" // <author>/<name>
}

func NewHandle(s string) (*Handle, error) {
	d := new(Handle)

	err := d.SetString(s)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// order: rsplit @, split /, rsplit .
func (d *Handle) SetString(s string) error {

	nam_idx := strings.Index(s, "/")
	if nam_idx < 0 {
		return HandleError(s, "no author/name separator")
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

	d.Handle = s

	// parts
	d.Author = slice(s, 0, nam_idx)
	d.Name = slice(s, nam_idx, fmt_idx)
	d.Format = slice(s, fmt_idx+1, ver_idx)
	d.Version = slice(s, ver_idx+1, len(s))
	d.Path = slice(s, 0, fmt_idx)
	return nil
}

func slice(s string, from int, to int) string {
	from = MaxInt(from, 0)
	to = MinInt(to, len(s))
	return s[MinInt(from, to):to]
}

// https://groups.google.com/forum/#!topic/golang-nuts/dbyqx_LGUxM is silly.
func MinInt(x, y int) (r int) {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x, y int) (r int) {
	if x > y {
		return x
	}
	return y
}

func (d *Handle) GoString() string {
	return d.Handle
}

func HandleError(handle string, problem string) error {
	return fmt.Errorf("Invalid handle (%s): %s", problem, handle)
}
