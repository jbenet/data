package data

import (
	"os"
	"regexp"
	"strings"
)

type InputField struct {
	Prompt  string
	Value   *string
	Pattern *regexp.Regexp
	Help    string
}

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
	fields := []InputField{
		InputField{
			"author id (required)",
			&h.Author,
			UserRegexp,
			"Must be a valid username. Can only contain [a-z0-9-_.].",
		},
		InputField{
			"dataset id (required)",
			&h.Name,
			IdentRegexp,
			"Must be a valid dataset id. Can only contain [a-z0-9-_.].",
		},
		InputField{
			"dataset version (required)",
			&h.Version,
			IdentRegexp,
			"Must be a valid version. Can only contain [a-z0-9-_.].",
		},
		InputField{"tagline description (required)", &df.Tagline, nil, ""},
		InputField{"long description (optional)", &df.Description, nil, ""},
		InputField{"license name (optional)", &df.License, nil, ""},
	}

	for _, field := range fields {
		err := fillOutField(field)
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

func fillOutField(f InputField) error {

	// validator function
	valid := func(val string) bool {
		if strings.Contains(f.Prompt, "required") && len(val) < 1 {
			return false
		}

		if f.Pattern != nil && !f.Pattern.MatchString(val) {
			return false
		}

		return true
	}

	for {
		pOut("Enter %s [%s]: ", f.Prompt, *f.Value)
		line, err := readInput()
		if err != nil {
			return err
		}

		// if not required, and entered nothing, get out.
		if len(line) == 0 && valid(*f.Value) {
			break
		}

		// if valid input
		if valid(line) {
			*f.Value = line
			break
		}

		if len(f.Help) > 0 {
			pOut("	Error: %s\n", f.Help)
		} else {
			pOut("	Error: Invalid input.\n")
		}
	}

	dOut("entered: %s\n", *f.Value)
	return nil
}

func fillOutUserProfile(p *UserProfile) error {
	pOut("Editing user profile. [Current value].\n")

	fields := []InputField{
		InputField{"Full Name", &p.Name, nil, ""},
		// "Email (required)":            &p.Email,
		InputField{"Website Url", &p.Website, nil, ""},
		InputField{"Github username", &p.Github, nil, ""},
		InputField{"Twitter username", &p.Twitter, nil, ""},
	}

	for _, f := range fields {
		err := fillOutField(f)
		if err != nil {
			return err
		}
	}

	return nil
}
