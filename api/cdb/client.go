package cdb

import (
	"encoding/json"
	"fmt"
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

func (c *CDBClient) get(suffix string, data any, queryparams map[string]string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.opts.BaseURL, suffix), nil)
	if err != nil {
		return err
	}

	if c.opts.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.opts.Token))
	}

	params := url.Values{}
	for k, v := range queryparams {
		params.Add(k, v)
	}
	req.URL.RawQuery = params.Encode()

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

	params := make(map[string]string)
	if q.Type != "" {
		params["type"] = string(q.Type)
	}
	if q.Query != "" {
		params["q"] = q.Query
	}

	err := c.get("api/packages", &pkgs, params)
	return pkgs, err
}

func (c *CDBClient) GetDetails(p *Package) (*PackageDetails, error) {
	details := &PackageDetails{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s", p.Author, p.Name), details, nil)
	return details, err
}

func (c *CDBClient) GetDependencies(p *Package) (PackageDependency, error) {
	pd := PackageDependency{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s/dependencies", p.Author, p.Name), &pd, nil)
	return pd, err
}

func (c *CDBClient) GetReleases(p *Package) ([]*PackageRelease, error) {
	pr := []*PackageRelease{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s/releases", p.Author, p.Name), &pr, nil)
	return pr, err
}

func (c *CDBClient) GetScreenshots(p *Package) ([]*PackageScreenshot, error) {
	ps := []*PackageScreenshot{}
	err := c.get(fmt.Sprintf("api/packages/%s/%s/screenshots", p.Author, p.Name), &ps, nil)
	return ps, err
}

func (c *CDBClient) GetThumbnails(ss *PackageScreenshot) *PackageThumbnails {
	return &PackageThumbnails{
		Small:  strings.ReplaceAll(ss.URL, "/uploads/", "/thumbnails/1/"),
		Medium: strings.ReplaceAll(ss.URL, "/uploads/", "/thumbnails/2/"),
		Large:  strings.ReplaceAll(ss.URL, "/uploads/", "/thumbnails/3/"),
	}
}
