package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient()
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

type mockAuth struct{}

func (mockAuth) ApplyAuthentication(r *http.Request) {
	r.Header.Add("Authentication", "applied")
}

// Test Client.newRequest when supplying JSON-serializable object for request body
func TestNewRequestJSON(t *testing.T) {
	c := NewClient()

	expectedMethod := "GET"
	inURL, expectedURL := "foo", defaultBaseURL+"foo"
	inBody, expectedBody := map[string]string{"foo": "bar"}, "{\"foo\":\"bar\"}\n"

	r, err := c.newRequest(expectedMethod, inURL, mockAuth{}, inBody)

	body, _ := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.Equal(t, expectedMethod, r.Method)
	assert.Equal(t, expectedURL, r.URL.String())
	assert.Equal(t, expectedBody, string(body))
	assert.Equal(t, "applied", r.Header.Get("Authentication"), "Authentication not applied")
}

// Test Client.newRequest when supplying io.Reader for request body
func TestNewRequestReader(t *testing.T) {
	c := NewClient()

	expectedMethod := "GET"
	inURL, expectedURL := "foo", defaultBaseURL+"foo"
	inBody, expectedBody := bytes.NewBufferString("foobar"), "foobar"

	r, err := c.newRequest(expectedMethod, inURL, mockAuth{}, inBody)

	body, _ := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.Equal(t, expectedMethod, r.Method)
	assert.Equal(t, expectedURL, r.URL.String())
	assert.Equal(t, expectedBody, string(body))
	assert.Equal(t, "applied", r.Header.Get("Authentication"), "Authentication not applied")
}

// Test Client.do when supplying an object for JSON deserialization of response body
func TestDoJSON(t *testing.T) {
	setup()
	defer teardown()

	expected := map[string]string{"foo": "bar"}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(expected)
	})

	req, _ := client.newRequest("GET", "/", NilAuth{}, nil)
	body := map[string]string{}

	_, err := client.do(req, &body)

	assert.NoError(t, err)
	assert.Equal(t, expected, body)
}

// Test Client.do when supplying an io.Writer for receipt of response body
func TestDoWriter(t *testing.T) {
	setup()
	defer teardown()

	str := "foo_bar"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, str)
	})

	req, _ := client.newRequest("GET", "/", NilAuth{}, nil)
	body := bytes.Buffer{}

	_, err := client.do(req, &body)

	assert.NoError(t, err)
	assert.Equal(t, str, body.String())
}
