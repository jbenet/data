package data

import (
	"path"
)

/*
  # Datafile format
  # A YAML (inc json) doc with the following keys:

  # required
  handle: <author>/<name>[.<format>][@<tag>]
  title: Dataset Title

  # optional functionality
  dependencies: [<other dataset handles>]
  formats: {<format> : <format url>}

  # optional information
  description: Text describing dataset.
  repository: <repo url>
  website: <dataset url>
  license: <license url>
  contributors: ["Author Name [<email>] [(url)]>", ...]
  sources: [<source urls>]
*/

// Serializable into YAML
type datafileContents struct {
	Dataset string
	Title   string ",omitempty"

	Mirrors      []string          ",omitempty"
	Dependencies []string          ",omitempty"
	Formats      map[string]string ",omitempty"

	Description  string   ",omitempty"
	Repository   string   ",omitempty"
	Website      string   ",omitempty"
	License      string   ",omitempty"
	Contributors []string ",omitempty"
	Sources      []string ",omitempty"
}

type Datafile struct {
	SerializedFile   "-"
	datafileContents ",inline"
}

const DatasetDir = "datasets"
const DatafileName = "Datafile"

func DatafilePath(dataset string) string {
	return path.Join(DatasetDir, dataset, DatafileName)
}

func NewDatafile(path string) (*Datafile, error) {
	df := &Datafile{SerializedFile: SerializedFile{Path: path}}
	df.SerializedFile.Format = df

	if len(path) > 0 {
		err := df.ReadFile()
		if err != nil {
			return df, err
		}
	}
	return df, nil
}

func NewDefaultDatafile() (*Datafile, error) {
	return NewDatafile(DatafileName)
}

func NewDatafileWithRef(ref string) (*Datafile, error) {
	f, _ := NewDatafile("")
	err := f.ReadBlob(ref)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (d *Datafile) Handle() *Handle {
	return NewHandle(d.Dataset)
}

func (d *Datafile) Valid() bool {
	return d.Handle().Valid()
}
