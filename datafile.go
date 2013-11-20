package main

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
type DatafileContents struct {
	Handle Handle ",inline"
	Title  string ",omitempty"

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
	path             string "-" // YAML ignore
	DatafileContents ",inline"
}

const DatasetDir = "datasets"

func NewDatafile(dataset string) (*Datafile, error) {
	df := &Datafile{path: path.Join(DatasetDir, dataset, "Datafile")}
	err := df.ReadFile()
	if err != nil {
		return nil, err
	}
	return df, nil
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

	// need to ensure handle is initialized...
	d.Handle.SetString(d.Handle.Handle)
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

	return ioutil.WriteFile(d.path, buf, 0666)
}

func (d *Datafile) ReadFile() error {
	buf, err := ioutil.ReadFile(d.path)
	if err != nil {
		return err
	}

	return d.Unmarshal(buf)
}
