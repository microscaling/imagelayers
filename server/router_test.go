package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	srv *httptest.Server
)

func mockHandler(w http.ResponseWriter, r *http.Request) {
}

func setupServer() {
	router := Router{mux.NewRouter()}
	m := RouteMap{
		"GET": {
			"/it": mockHandler,
		},
	}

	router.AddRoutes("/test", m)
	router.AddCorsRoutes("/cors", m)
	srv = httptest.NewServer(router)
}

func teardownServer() {
	srv.Close()
}

func TestRoute(t *testing.T) {
	setupServer()
	defer teardownServer()

	path := srv.URL + "/test/it"
	resp, err := http.DefaultClient.Get(path)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCorsRoute(t *testing.T) {
	setupServer()
	defer teardownServer()

	path := srv.URL + "/cors/it"

	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Add("Origin", srv.URL)
	resp, err := http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, srv.URL, resp.Header.Get("Access-Control-Allow-Origin"), "CORS Allow-Origin not set")
}
