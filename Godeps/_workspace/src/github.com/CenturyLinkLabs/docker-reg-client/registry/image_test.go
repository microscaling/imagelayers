package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestImageGetLayer(t *testing.T) {
	setup()
	defer teardown()

	expected := "_foo_\n_bar_"

	mux.HandleFunc("/images/1234567890/layer", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, expected)
	})

	buffer := &bytes.Buffer{}
	err := client.Image.GetLayer("1234567890", buffer, NilAuth{})

	assert.NoError(t, err)
	assert.Equal(t, expected, buffer.String())
}

func TestImageAddLayer(t *testing.T) {
	setup()
	defer teardown()

	expected := "_foo_\n_bar_"

	mux.HandleFunc("/images/1234567890/layer", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)

		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expected, string(body))
	})

	buffer := bytes.NewBufferString(expected)
	err := client.Image.AddLayer("1234567890", buffer, NilAuth{})

	assert.NoError(t, err)
}

func TestImageGetMetadata(t *testing.T) {
	setup()
	defer teardown()

	loc, _ := time.LoadLocation("")
	expected := &ImageMetadata{
		ID:      "1234567890",
		Parent:  "0987654321",
		Comment: "lgtm",
		Created: time.Date(1973, 10, 25, 8, 0, 0, 0, loc),
		Author:  "Brian DeHamer",
		Size:    20909309,
	}

	mux.HandleFunc("/images/1234567890/json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(expected)
	})

	meta, err := client.Image.GetMetadata("1234567890", NilAuth{})

	assert.NoError(t, err)
	assert.Equal(t, expected, meta)
}

func TestImageAddMetadata(t *testing.T) {
	setup()
	defer teardown()

	meta := &ImageMetadata{
		ID:      "1234567890",
		Parent:  "0987654321",
		Comment: "lgtm",
		Created: time.Now(),
		Author:  "Brian DeHamer",
		Size:    20909309,
	}

	mux.HandleFunc("/images/1234567890/json", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)

		expected, _ := json.Marshal(meta)
		expected = append(expected, 0xa)
		assert.Equal(t, expected, body)
	})

	err := client.Image.AddMetadata("1234567890", meta, NilAuth{})

	assert.NoError(t, err)
}

func TestImageGetAncestry(t *testing.T) {
	setup()
	defer teardown()

	expected := []string{"12345", "67890"}

	mux.HandleFunc("/images/1234567890/ancestry", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(expected)
	})

	ancestry, err := client.Image.GetAncestry("1234567890", NilAuth{})

	assert.NoError(t, err)
	assert.Equal(t, expected, ancestry)
}
