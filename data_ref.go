package data

import (
	"fmt"
	"io/ioutil"
)

const RefLatest = "latest"

// serializable into YAML
type DatasetRefs struct {

	// All published refs are listed here. { ref-hash : iso-timestamp }
	Published map[string]string

	// Automatic named pointers to published references. { version : ref-hash }
	// Generated from dataset handle versions.
	Versions map[string]string
}

func (r DatasetRefs) LastUpdated() string {
	pl := sortMapByValue(r.Published)
	if len(pl) > 0 {
		return pl[len(pl)-1].Value
	}
	return ""
}

func (r DatasetRefs) LatestPublished() string {
	s := r.SortedPublished()
	if len(s) == 0 {
		return ""
	}
	return s[len(s)-1]
}

func (r DatasetRefs) SortedPublished() []string {
	vs := []string{}
	pl := sortMapByValue(r.Published)
	for _, p := range pl {
		vs = append(vs, p.Key)
	}
	return vs
}

func (r DatasetRefs) ResolveRef(ref string) string {

	// default to latest (like HEAD)
	if len(ref) == 0 {
		ref = RefLatest
	}

	// latest -> timestamp sorted
	if ref == RefLatest {
		return r.LatestPublished()
	}

	// look it up in versions table
	ref2, found := r.Versions[ref]
	if found {
		return ref2
	}

	// Guess we have no link, return it then.
	return ref
}

type HttpRefIndex struct {
	Http    *HttpClient
	Dataset string
	Refs    *DatasetRefs
}

func (h *HttpRefIndex) FetchRefs(refresh bool) error {
	// already fetched?
	if h.Refs != nil && !refresh {
		return nil
	}

	resp, err := h.Http.Get("")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	refs := &DatasetRefs{}
	err = Unmarshal(resp.Body, refs)
	if err != nil {
		return err
	}

	// set at the end, once we're sure no errors happened
	h.Refs = refs
	return nil
}

func (h *HttpRefIndex) Has(ref string) (bool, error) {
	return httpExists(h.Http.SubUrl(ref))
}

func (h *HttpRefIndex) Get(ref string) (string, error) {
	resp, err := h.Http.Get(ref)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf[:]), nil
}

func (h *HttpRefIndex) Put(ref string) error {
	resp, err := h.Http.Post(ref, nil)
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

func (h *HttpRefIndex) VersionRef(version string) (string, error) {
	err := h.FetchRefs(false)
	if err != nil {
		return "", err
	}

	// special ref latest
	if version == RefLatest {
		refs := h.SortedPublished()
		if len(refs) == 0 {
			return "", fmt.Errorf("no published refs")
		}
		return refs[len(refs)-1], nil
	}

	ref, _ := h.Refs.Versions[version]
	return ref, nil
}

func (h *HttpRefIndex) RefTimestamp(ref string) (string, error) {
	err := h.FetchRefs(false)
	if err != nil {
		return "", err
	}

	time, _ := h.Refs.Published[ref]
	return time, nil
}

func (h *HttpRefIndex) SortedPublished() []string {
	return h.Refs.SortedPublished()
}

// DataIndex extension to generate a RefIndex
func (d *DataIndex) RefIndex(dataset string) *HttpRefIndex {
	ri := &HttpRefIndex{
		Http: &HttpClient{
			Url:       d.Http.Url + "/" + dataset + "/" + "refs",
			User:      d.Http.User,
			AuthToken: d.Http.AuthToken,
		},
		Dataset: dataset,
	}
	return ri
}
