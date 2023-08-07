package cdb

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *CDBClient) GetPackages() ([]*Package, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/packages", c.opts.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	if c.opts.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.opts.Token))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	p := make([]*Package, 0)
	err = json.NewDecoder(resp.Body).Decode(&p)

	return p, err
}
