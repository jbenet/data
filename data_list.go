package main

import (
	"io/ioutil"
	"path"
)

const DatasetDir = "datasets"

func ListCmd() {
	ListDatasets(DatasetDir)
}

func ListDatasets(dir string) error {
	authors, err := ioutil.ReadDir(dir)

	if err != nil {
		DErr("data: error reading dataset directory \"%s\"\n", dir)
		return err
	}

	for _, a := range authors {

		author := path.Join(dir, a.Name())
		datasets, err := ioutil.ReadDir(author)
		if err != nil {
			continue
		}

		for _, d := range datasets {
			dataset := path.Join(author, d.Name())
			datafile := &Datafile{path: path.Join(dataset, "Datafile")}

			err := datafile.ReadFile()
			if err != nil {
				DErr("%s\n", err)
				continue
			}

			DOut("%s\n", datafile.Handle)
		}
	}

	return nil
}
