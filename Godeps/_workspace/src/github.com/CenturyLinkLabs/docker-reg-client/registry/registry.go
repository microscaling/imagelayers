package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://index.docker.io/v1/"
)

// Client manages communication with the Docker Registry API.
type Client struct {
	// BaseURL represents the base URL for the Registry API.
	// Defaults to "https://index.docker.io/v1/".
	BaseURL *url.URL

	// Hub gives access to the Docker Hub API for retrieving auth tokens.
	Hub *HubService

	// Image gives access to the /images part of the Registry API.
	Image *ImageService

	// Repository gives access to the /repositories part of the Registry API.
	Repository *RepositoryService

	// Search gives access to the /search part of the Registry API.
	Search *SearchService

	client *http.Client
}

// RegistryError encapsulates any errors which result from communicating with
// the Docker Registry API
type RegistryError struct {
	Code   int
	method string
	url    *url.URL
}

func (e RegistryError) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.method, e.url, e.Code)
}

// NewClient returns a new Docker Registry API client.
func NewClient() *Client {
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		panic(err)
	}

	c := &Client{
		BaseURL: baseURL,
		client:  http.DefaultClient,
	}

	c.Hub = &HubService{client: c}
	c.Image = &ImageService{client: c}
	c.Repository = &RepositoryService{client: c}
	c.Search = &SearchService{client: c}

	return c
}

// Creates a new API request. Relative URLs should be specified without a
// preceding slash.
func (c *Client) newRequest(method, urlStr string, auth Authenticator, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(rel)

	buf := &bytes.Buffer{}
	if body != nil {
		if reader, ok := body.(io.Reader); ok {
			io.Copy(buf, reader)
		} else {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	auth.ApplyAuthentication(req)

	return req, nil
}

// Performs the given request and decodes the response into the variable v.
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if status := resp.StatusCode; status < 200 || status > 299 {
		return resp, RegistryError{Code: status, method: req.Method, url: req.URL}
	}

	if v != nil {
		if writer, ok := v.(io.Writer); ok {
			io.Copy(writer, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}
