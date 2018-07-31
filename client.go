package synapse

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
}

func newClient(URL string) (*Client, error) {
	url, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	return &Client{
		URL:        url,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *Client) Match(ctx context.Context, name string, params Context) (Hits, error) {
	j, err := json.Marshal(params)
	if err != nil {
		return Hits{}, err
	}

	req, err := c.newRequest(ctx, "GET", "", bytes.NewReader(j))
	if err != nil {
		return Hits{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return Hits{}, err
	}

	var hits Hits
	if err := c.decodeBody(res, &hits); err != nil {
		return Hits{}, err
	}

	return hits, nil
}

func (c *Client) newRequest(ctx context.Context, method, urlPath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, urlPath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	return req.WithContext(ctx), nil
}

func (c Client) decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
