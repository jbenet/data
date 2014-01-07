package data

import (
	"strings"
)

// serializable into YAML
type DatasetRefs struct {

	// All published refs are listed here. { ref-hash : iso-timestamp }
	Published map[string]string

	// Automatic named pointers to published references. { version : ref-hash }
	// Generated from dataset handle versions.
	Versions map[string]string
}

type HttpRefIndex struct {
	Dataset string
	BaseUrl string
	Refs    *DatasetRefs
}

func (h HttpRefIndex) Url(url string) string {
	return h.BaseUrl + "/" + h.Dataset + "/refs/" + url
}

func (h *HttpRefIndex) FetchRefs(refresh bool) error {
	// already fetched?
	if h.Refs != nil && !refresh {
		return nil
	}

	resp, err := httpGet(h.Url(""))
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
	return httpExists(h.Url(ref))
}

func (h *HttpRefIndex) Get(ref string) (string, error) {
	buf, err := httpReadAll(h.Url(ref))
	if err != nil {
		return "", err
	}

	return string(buf[:]), nil
}

func (h *HttpRefIndex) Put(ref string) error {
	resp, err := httpPost(h.Url(ref), "text", strings.NewReader(""))
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

// DataIndex extension to generate a RefIndex
func (d *DataIndex) RefIndex(dataset string) *HttpRefIndex {
	return &HttpRefIndex{
		Dataset: dataset,
		BaseUrl: d.Url,
	}
}
