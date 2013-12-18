package data

import (
	"github.com/jbenet/commander"
	"io/ioutil"
	"path"
)

var cmd_data_list = &commander.Command{
	UsageLine: "info <dataset>",
	Short:     "Show dataset information.",
	Long: `data info - Show dataset information.

    Returns the Datafile corresponding to <dataset> and exits.
  `,
	Run: listCmd,
}

func listCmd(*commander.Command, []string) error {
	return listDatasets(DatasetDir)
}

func listDatasets(dir string) error {
	authors, err := ioutil.ReadDir(dir)

	if err != nil {
		pErr("data: error reading dataset directory \"%s\"\n", dir)
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
				pErr("Error: %s\n", err)
				continue
			}

			h := datafile.Handle()
			pOut("    %-20s @%s\n", h.Path(), h.Version)
		}
	}

	return nil
}
