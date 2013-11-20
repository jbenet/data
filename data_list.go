package main

import (
	"io/ioutil"
	"path"
)

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
			dataset := path.Join(a.Name(), d.Name())
			datafile, err := NewDatafile(dataset)
			if err != nil {
				DErr("Error: %s\n", err)
				continue
			}

			DOut("    %-20s @%s\n", datafile.Handle.Path, datafile.Handle.Version)
		}
	}

	return nil
}
