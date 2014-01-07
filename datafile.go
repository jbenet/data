package data

import (
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

// Serializable into YAML
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
	SerializedFile   "-"
	datafileContents ",inline"
}

const DatasetDir = "datasets"
const DatasetFile = "Datafile"

func DatafilePath(dataset string) string {
	return path.Join(DatasetDir, dataset, DatasetFile)
}

func NewDatafile(path string) (*Datafile, error) {
	if len(path) < 1 {
		path = DatasetFile
	}

	df := &Datafile{SerializedFile: SerializedFile{Path: path}}
	df.SerializedFile.Format = df

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
