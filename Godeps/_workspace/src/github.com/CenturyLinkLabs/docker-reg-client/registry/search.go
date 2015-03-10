package registry

import (
	"fmt"
	"net/url"
	"strconv"
)

// SearchService gives access to the /search portion of the Registry API
type SearchService struct {
	client *Client
}

type SearchResults struct {
	NumPages   int            `json:"num_pages"`
	NumResults int            `json:"num_results"`
	Results    []SearchResult `json:"results"`
	PageSize   int            `json:"page_size"`
	Query      string         `json:"query"`
	Page       int            `json:"page,string"`
}

type SearchResult struct {
	Name        string `json:"name"`
	IsAutomated bool   `json:"is_automated"`
	IsTrusted   bool   `json:"is_trusted"`
	IsOfficial  bool   `json:"is_official"`
	StarCount   int    `json:"star_count"`
	Description string `json:"description"`
}

func (s *SearchService) Query(query string, page, num int) (*SearchResults, error) {
	params := url.Values{}
	params.Add("q", query)
	params.Add("page", strconv.Itoa(page))
	params.Add("n", strconv.Itoa(num))
	path := fmt.Sprintf("search?%s", params.Encode())

	req, err := s.client.newRequest("GET", path, NilAuth{}, nil)
	if err != nil {
		return nil, err
	}

	results := &SearchResults{}
	_, err = s.client.do(req, results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
