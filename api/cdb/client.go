package cdb

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CDBClient struct {
	client http.Client
	opts   *CDBClientOpts
}

type CDBClientOpts struct {
	BaseURL string
	Token   string
}

func NewWithOpts(opts *CDBClientOpts) *CDBClient {
	return &CDBClient{
		client: http.Client{},
		opts:   opts,
	}
}
func New() *CDBClient {
	return NewWithOpts(&CDBClientOpts{
		BaseURL: "https://content.minetest.net",
		Token:   "",
	})
}

func (c *CDBClient) get(suffix string, data any, params url.Values) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.opts.BaseURL, suffix), nil)
	if err != nil {
		return err
	}

	if c.opts.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.opts.Token))
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(data)
}

func (c *CDBClient) GetPackages() ([]*Package, error) {
	pkgs := make([]*Package, 0)
	err := c.get("api/packages", &pkgs, nil)
	return pkgs, err
}

func (c *CDBClient) SearchPackages(q *PackageQuery) ([]*Package, error) {
	pkgs := make([]*Package, 0)

	params := url.Values{}
	for _, t := range q.Type {
		params.Add("type", string(t))
	}
	if q.Query != "" {
		params.Add("q", q.Query)
	}
	if q.Author != "" {
		params.Add("author", q.Author)
	}
	if q.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", q.Limit))
	}
	for _, cw := range q.Hide {
		params.Add("hide", cw.Name)
	}
	if q.Sort != "" {
		params.Add("sort", string(q.Sort))
	}
	if q.Order != "" {
		params.Add("order", string(q.Order))
	}
	if q.ProtocolVersion > 0 {
		params.Add("protocol_version", fmt.Sprintf("%d", q.ProtocolVersion))
	}
	if q.EngineVersion != "" {
		params.Add("engine_version", q.EngineVersion)
	}

	err := c.get("api/packages", &pkgs, params)
	return pkgs, err
}

func (c *CDBClient) GetTags() ([]*Tag, error) {
	tags := make([]*Tag, 0)
	err := c.get("api/tags", &tags, nil)
	return tags, err
}

func (c *CDBClient) GetContentWarnings() ([]*ContentWarning, error) {
	list := make([]*ContentWarning, 0)
	err := c.get("api/content_warnings", &list, nil)
	return list, err
}

func (c *CDBClient) GetDetails(author, name string) (*PackageDetails, error) {
	details := &PackageDetails{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s", author, name), details, nil)
	return details, err
}

func (c *CDBClient) GetDependencies(author, name string) (PackageDependency, error) {
	pd := PackageDependency{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s/dependencies", author, name), &pd, nil)
	return pd, err
}

func (c *CDBClient) GetReleases(author, name string) ([]*PackageRelease, error) {
	pr := []*PackageRelease{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s/releases", author, name), &pr, nil)
	return pr, err
}

func (c *CDBClient) GetRelease(author, name string, id int) (*PackageRelease, error) {
	pr := &PackageRelease{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s/releases/%d", author, name, id), pr, nil)
	return pr, err
}

func (c *CDBClient) GetScreenshots(author, name string) ([]*PackageScreenshot, error) {
	ps := []*PackageScreenshot{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s/screenshots", author, name), &ps, nil)
	return ps, err
}

func (c *CDBClient) GetThumbnails(ss *PackageScreenshot) *PackageThumbnails {
	return &PackageThumbnails{
		Small:  strings.ReplaceAll(ss.URL, "/uploads/", "/thumbnails/1/"),
		Medium: strings.ReplaceAll(ss.URL, "/uploads/", "/thumbnails/2/"),
		Large:  strings.ReplaceAll(ss.URL, "/uploads/", "/thumbnails/3/"),
	}
}

func (c *CDBClient) GetDownloadURL(r *PackageRelease) string {
	return fmt.Sprintf("%s%s", c.opts.BaseURL, r.URL)
}

func (c *CDBClient) DownloadZip(r *PackageRelease) (*zip.Reader, error) {
	url := c.GetDownloadURL(r)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not download zip file: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not copy contents: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	return zip.NewReader(reader, reader.Size())
}
