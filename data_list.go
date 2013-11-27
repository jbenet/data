package data

import (
	"io/ioutil"
	"path"
)

func ListCmd([]string) error {
	return ListDatasets(DatasetDir)
}

func ListDatasets(dir string) error {
	authors, err := ioutil.ReadDir(dir)

	if err != nil {
		Err("data: error reading dataset directory \"%s\"\n", dir)
		return err
	}

	// for each author dir
	for _, a := range authors {
		// skip hidden files
		if a.Name()[0] == '.' {
			continue
		}

		author := path.Join(dir, a.Name())
		datasets, err := ioutil.ReadDir(author)
		if err != nil {
			continue
		}

		// for each dataset dir
		for _, d := range datasets {
			// skip hidden files
			if d.Name()[0] == '.' {
				continue
			}

			dataset := path.Join(a.Name(), d.Name())
			datafile, err := NewDatafile(DatafilePath(dataset))
			if err != nil {
				Err("Error: %s\n", err)
				continue
			}

			Out("    %-20s @%s\n", datafile.Handle.Path, datafile.Handle.Version)
		}
	}

	return nil
}
