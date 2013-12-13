package data

import (
	"bufio"
	"os"
	"strings"
)

func ensureDatafileInPath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	// if it doesn't exist, create it.
	f, err := os.Create(path)
	defer f.Close()

	return nil
}

func fillOutDatafileInPath(path string) error {

	err := ensureDatafileInPath(path)
	if err != nil {
		return err
	}

	df, err := NewDatafile(path)
	if err != nil {
		return err
	}

	return fillOutDatafile(df)
}

func fillOutDatafile(df *Datafile) error {
	pOut("Verifying Datafile fields...\n")

	h := df.Handle()
	fields := map[string]*string{
		"author name (required)":     &h.Author,
		"dataset id (required)":      &h.Name,
		"dataset version (required)": &h.Version,
		"dataset title (optional)":   &df.Title,
		"description (optional)":     &df.Description,
		"license name (optional)":    &df.License,
	}

	for p, f := range fields {
		err := fillOutDatafileField(p, f)
		if err != nil {
			return err
		}

		df.Dataset = h.Dataset()
		if df.Valid() {
			err = df.WriteFile()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fillOutDatafileField(prompt string, field *string) error {
	for len(*field) < 1 {
		pOut("Enter %s: ", prompt)
		line, err := readInput()
		if err != nil {
			return err
		}

		*field = line

		// if not required, don't loop
		if !strings.Contains(prompt, "required") {
			break
		}
	}

	dOut("entered: %s\n", *field)

	return nil
}

func readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(line), nil
}