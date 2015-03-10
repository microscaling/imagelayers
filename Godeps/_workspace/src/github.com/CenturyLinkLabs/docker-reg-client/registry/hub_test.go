package registry

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHubGetReadToken(t *testing.T) {
	setup()
	defer teardown()

	token := "foo123"
	host := "www.bar.com"

	mux.HandleFunc("/repositories/foo/images", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "true", r.Header.Get("X-Docker-Token"))

		w.Header().Add("X-Docker-Token", token)
		w.Header().Add("X-Docker-Endpoints", host)
	})

	auth, err := client.Hub.GetReadToken("foo")

	assert.NoError(t, err)
	assert.Equal(t, Read, auth.Access)
	assert.Equal(t, token, auth.Token)
	assert.Equal(t, host, auth.Host)
}

func TestHubGetWriteToken(t *testing.T) {
	setup()
	defer teardown()

	token := "foo123"
	host := "www.bar.com"

	mux.HandleFunc("/repositories/foo/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "true", r.Header.Get("X-Docker-Token"))

		w.Header().Add("X-Docker-Token", token)
		w.Header().Add("X-Docker-Endpoints", host)
	})

	auth, err := client.Hub.GetWriteToken("foo", NilAuth{})

	assert.NoError(t, err)
	assert.Equal(t, Write, auth.Access)
	assert.Equal(t, token, auth.Token)
	assert.Equal(t, host, auth.Host)
}

func TestHubGetDeleteToken(t *testing.T) {
	setup()
	defer teardown()

	token := "foo123"
	host := "www.bar.com"

	mux.HandleFunc("/repositories/foo/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "true", r.Header.Get("X-Docker-Token"))

		w.Header().Add("X-Docker-Token", token)
		w.Header().Add("X-Docker-Endpoints", host)
	})

	auth, err := client.Hub.GetDeleteToken("foo", NilAuth{})

	assert.NoError(t, err)
	assert.Equal(t, Delete, auth.Access)
	assert.Equal(t, token, auth.Token)
	assert.Equal(t, host, auth.Host)
}
