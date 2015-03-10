package registry

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	username := "foo"
	password := "bar"
	auth := BasicAuth{username, password}

	req, _ := http.NewRequest("GET", "http://foo.com/v1", nil)
	auth.ApplyAuthentication(req)

	u, p, _ := req.BasicAuth()
	assert.Equal(t, username, u)
	assert.Equal(t, password, p)
}

func TestTokenAuth(t *testing.T) {
	host := "bar.com"
	token := "token xyz"
	auth := TokenAuth{Host: host, Token: token}

	req, _ := http.NewRequest("GET", "http://foo.com/v1", nil)
	auth.ApplyAuthentication(req)

	assert.Equal(t, "Token "+token, req.Header.Get("Authorization"))
	assert.Equal(t, host, req.Host)
	assert.Equal(t, host, req.URL.Host)
}
