package registry

import (
	"fmt"
	"net/http"
	"net/url"
)

// HubService gives acess to the Docker Hub API used to retrieve authentication tokens.
type HubService struct {
	client *Client
}

// Retrieves a read-only token from the Docker Hub for the specified repo.
func (h *HubService) GetReadToken(repo string) (*TokenAuth, error) {
	path := fmt.Sprintf("repositories/%s/images", repo)
	req, err := h.newRequest("GET", path, NilAuth{})
	if err != nil {
		return nil, err
	}

	return h.do(req)
}

// Retrieves a read-only token from the Docker Hub for the specified private repo.
func (h *HubService) GetReadTokenWithAuth(repo string, auth Authenticator) (*TokenAuth, error) {
	path := fmt.Sprintf("repositories/%s/images", repo)
	req, err := h.newRequest("GET", path, auth)
	if err != nil {
		return nil, err
	}

	return h.do(req)
}

// Retrieves a write-only token from the Docker Hub for the specified repo. The
// auth argument should be the user's basic auth credentials.
func (h *HubService) GetWriteToken(repo string, auth Authenticator) (*TokenAuth, error) {
	path := fmt.Sprintf("repositories/%s/", repo)
	req, err := h.newRequest("PUT", path, auth)
	if err != nil {
		return nil, err
	}

	return h.do(req)
}

// Retrieves a write-only token from the Docker Hub for the specified repo. The
// auth argument should be the user's basic auth credentials.
func (h *HubService) GetDeleteToken(repo string, auth Authenticator) (*TokenAuth, error) {
	path := fmt.Sprintf("repositories/%s/", repo)
	req, err := h.newRequest("DELETE", path, auth)
	if err != nil {
		return nil, err
	}

	return h.do(req)
}

func (h *HubService) newRequest(method, urlStr string, auth Authenticator) (*http.Request, error) {
	req, err := h.client.newRequest(method, urlStr, auth, []string{})
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Docker-Token", "true")

	return req, nil
}

func (h *HubService) do(req *http.Request) (*TokenAuth, error) {
	res, err := h.client.do(req, nil)
	if err != nil {
		return nil, err
	}

	token := &TokenAuth{Token: res.Header.Get("X-Docker-Token")}

	switch req.Method {
	case "GET":
		token.Access = Read
	case "PUT":
		token.Access = Write
	case "DELETE":
		token.Access = Delete
	}

	// Need to parse the value cause different requests return the endpoints
	// in different formats
	url, _ := url.Parse(res.Header.Get("X-Docker-Endpoints"))
	if url.Host != "" {
		token.Host = url.Host
	} else {
		token.Host = url.Path
	}

	return token, nil
}
