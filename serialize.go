package data

import (
	"io"
	"io/ioutil"
	"launchpad.net/goyaml"
)

type file struct {
	Path   string      "-"
	format interface{} "-"
}

func (f *file) Marshal() ([]byte, error) {
	dOut("Marshalling %s\n", f.Path)
	return goyaml.Marshal(f.format)
}

func (f *file) Unmarshal(buf []byte) error {
	err := goyaml.Unmarshal(buf, f.format)
	if err != nil {
		return err
	}

	dOut("Unmarshalling %s\n", f.Path)
	return nil
}

func (f *file) Write(w io.Writer) error {
	buf, err := f.Marshal()
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	return err
}

func (f *file) Read(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return f.Unmarshal(buf)
}

func (f *file) WriteFile() error {
	buf, err := f.Marshal()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.Path, buf, 0666)
}

func (f *file) ReadFile() error {
	buf, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}

	return f.Unmarshal(buf)
}
