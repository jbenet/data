package data

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type DataIndex struct {
	Http      *HttpClient
	BlobStore blobStore
}

var mainDataIndex *DataIndex

const mainIndexName = "datadex"

func (i *DataIndex) ArchiveUrl(h *Handle) string {
	ref := "master"
	if len(h.Version) > 0 {
		ref = h.Version
	}
	return i.Http.SubUrl(path.Join(h.Path(), "archive", ref, ArchiveSuffix))
}

// why not use `func init()`? some commands don't need an index
// is annoying to error out on an S3 key when S3 isn't needed.
func NewMainDataIndex() (*DataIndex, error) {
	if mainDataIndex != nil {
		return mainDataIndex, nil
	}

	blobStore, err := NewS3Store("datadex.archives")
	if err != nil {
		return nil, err
	}

	h, err := NewHttpClient()
	if err != nil {
		return nil, err
	}

	mainDataIndex := &DataIndex{Http: h}
	mainDataIndex.BlobStore = blobStore
	return mainDataIndex, nil
}

const HttpHeaderUser = "X-Data-User"
const HttpHeaderToken = "X-Data-Token"
const HttpHeaderContentType = "Content-Type"
const HttpHeaderContentTypeYaml = "application/yaml"

// Controls authenticated http accesses.
type HttpClient struct {
	Url       string
	User      string
	AuthToken string
}

func NewHttpClient() (*HttpClient, error) {
	i, exists := Config.Index[mainIndexName]
	if !exists {
		return nil, fmt.Errorf("Config error: no datadex index.")
	}

	h := &HttpClient{
		Url:       i.Url,
		User:      i.User,
		AuthToken: i.Token,
	}

	return h, nil
}

func (h HttpClient) SubUrl(path string) string {
	return h.Url + "/" + path
}

func (h *HttpClient) Get(path string) (*http.Response, error) {
	dOut("http index get %s\n", h.SubUrl(path))

	req, err := http.NewRequest("GET", h.SubUrl(path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(HttpHeaderToken, h.AuthToken)
	req.Header.Add(HttpHeaderUser, h.User)
	return h.DoRequest(req)
}

func (h *HttpClient) Post(path string, body interface{}) (*http.Response, error) {
	dOut("http index post %s\n", h.SubUrl(path))

	rdr := io.Reader(nil)
	var err error
	if body != nil {
		rdr, err = Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("POST", h.SubUrl(path), rdr)
	if err != nil {
		return nil, err
	}

	req.Header.Add(HttpHeaderContentType, HttpHeaderContentTypeYaml)
	req.Header.Add(HttpHeaderToken, h.AuthToken)
	req.Header.Add(HttpHeaderUser, h.User)
	return h.DoRequest(req)
}

func (h *HttpClient) DoRequest(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	c := resp.StatusCode
	if 200 <= c && c < 400 {
		return resp, nil
	}

	e, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	s := strings.TrimSpace(string(e[:]))
	return nil, fmt.Errorf("HTTP error status code: %d (%s)", c, s)
}
