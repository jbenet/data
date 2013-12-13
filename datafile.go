package data

import (
	"io"
	"io/ioutil"
	"launchpad.net/goyaml"
	"path"
)

/*
  Datafile format

  A YAML (inc json) doc with the following keys:

  (required:)
  handle: <author>/<name>[.<format>][@<tag>]
  title: Dataset Title

  (optional functionality:)
  dependencies: [<other dataset handles>]
  formats: {<format> : <format url>}

  (optional information:)
  description: Text describing dataset.
  repository: <repo url>
  homepage: <dataset url>
  license: <license url>
  contributors: ["Author Name [<email>] [(url)]>", ...]
  sources: [<source urls>]
*/

// Serializbale into YAML
type datafileContents struct {
	Dataset string
	Title   string ",omitempty"

	Mirrors      []string          ",omitempty"
	Dependencies []string          ",omitempty"
	Formats      map[string]string ",omitempty"

	Description  string   ",omitempty"
	Repository   string   ",omitempty"
	Homepage     string   ",omitempty"
	License      string   ",omitempty"
	Contributors []string ",omitempty"
	Sources      []string ",omitempty"
}

type Datafile struct {
	Path             string "-" // YAML ignore
	datafileContents ",inline"
}

const DatasetDir = "datasets"
const DatasetFile = "Datafile"

func DatafilePath(dataset string) string {
	return path.Join(DatasetDir, dataset, DatasetFile)
}

func NewDatafile(path string) (*Datafile, error) {
	df := &Datafile{Path: path}
	err := df.ReadFile()
	if err != nil {
		return nil, err
	}
	return df, nil
}

func (d *Datafile) Handle() *Handle {
	return NewHandle(d.Dataset)
}

func (d *Datafile) Valid() bool {
	return d.Handle().Valid()
}

// Serializing in/out

func (d *Datafile) Marshal() ([]byte, error) {
	return goyaml.Marshal(d)
}

func (d *Datafile) Unmarshal(buf []byte) error {
	err := goyaml.Unmarshal(buf, d)
	if err != nil {
		return err
	}

	return nil
}

func (d *Datafile) Write(w io.Writer) error {
	buf, err := d.Marshal()
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	return err
}

func (d *Datafile) Read(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return d.Unmarshal(buf)
}

func (d *Datafile) WriteFile() error {
	buf, err := d.Marshal()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(d.Path, buf, 0666)
}

func (d *Datafile) ReadFile() error {
	buf, err := ioutil.ReadFile(d.Path)
	if err != nil {
		return err
	}

	return d.Unmarshal(buf)
}
