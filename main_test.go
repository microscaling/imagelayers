package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	srv *httptest.Server
)

func setupServer() {
	srv = httptest.NewServer(NewServer().createRouter())
}

func teardownServer() {
	srv.Close()
}

func TestGraphRoute(t *testing.T) {
	setupServer()
	defer teardownServer()

	path := srv.URL
	resp, err := http.DefaultClient.Get(path)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}


func TestStatusGetSuccess (t *testing.T) {
	setupServer()
	defer teardownServer()

	path := srv.URL + "/registry/status"

	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Add("Origin", srv.URL)
	resp, err := http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, srv.URL, resp.Header.Get("Access-Control-Allow-Origin"), "CORS Allow-Origin not set")
}


func TestAnalyzePostSuccess (t *testing.T) {
	setupServer()
	defer teardownServer()

	path := srv.URL + "/registry/analyze"

	req, _ := http.NewRequest("POST", path, strings.NewReader("{}"))
	req.Header.Add("Origin", srv.URL)
	resp, err := http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, srv.URL, resp.Header.Get("Access-Control-Allow-Origin"), "CORS Allow-Origin not set")
}

func TestAnalyzeGetNotFound (t *testing.T) {
	setupServer()
	defer teardownServer()

	path := srv.URL + "/registry/analyze"

	req, _ := http.NewRequest("GET", path, nil)
	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestStatusPostNotFound (t *testing.T) {
	setupServer()
	defer teardownServer()

	path := srv.URL + "/registry/status"

	req, _ := http.NewRequest("POST", path, strings.NewReader("{}"))
	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestBadURL (t *testing.T) {
	setupServer()
	defer teardownServer()

	path := srv.URL + "/registry/foo"

	req, _ := http.NewRequest("GET", path, nil)
	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, resp.StatusCode)
}
