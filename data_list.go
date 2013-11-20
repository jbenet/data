package main

import (
	"io/ioutil"
	"path"
)

const DatasetDir = "datasets"

func ListCmd([]string) {
	ListDatasets(DatasetDir)
}

func ListDatasets(dir string) error {
	authors, err := ioutil.ReadDir(dir)

	if err != nil {
		DErr("data: error reading dataset directory \"%s\"\n", dir)
		return err
	}

	// for each author dir
	for _, a := range authors {
		author := path.Join(dir, a.Name())
		datasets, err := ioutil.ReadDir(author)
		if err != nil {
			continue
		}

		// for each dataset dir
		for _, d := range datasets {
			dataset := path.Join(author, d.Name())
			datafile := &Datafile{path: path.Join(dataset, "Datafile")}

			err := datafile.ReadFile()
			if err != nil {
				DErr("Error: %s\n", err)
				continue
			}

			DOut("    %-20s @%s\n", datafile.Handle.Path, datafile.Handle.Version)
		}
	}

	return nil
}
