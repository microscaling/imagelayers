package registry

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryListTags(t *testing.T) {
	setup()
	defer teardown()

	expected := TagMap{"bar": "1234567890"}

	mux.HandleFunc("/repositories/foo/tags", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(expected)
	})

	tags, err := client.Repository.ListTags("foo", NilAuth{})

	assert.NoError(t, err)
	assert.Equal(t, expected, tags)
}

func TestRepositoryGetImageID(t *testing.T) {
	setup()
	defer teardown()

	expected := "1234567890"

	mux.HandleFunc("/repositories/foo/tags/bar", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(expected)
	})

	imageID, err := client.Repository.GetImageID("foo", "bar", NilAuth{})

	assert.NoError(t, err)
	assert.Equal(t, expected, imageID)
}

func TestRepositorySetTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repositories/foo/tags/bar", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)

		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, "\"1234567890\"\n", string(body))
	})

	err := client.Repository.SetTag("foo", "1234567890", "bar", NilAuth{})

	assert.NoError(t, err)
}

func TestRepositoryDeleteTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repositories/foo/tags/bar", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
	})

	err := client.Repository.DeleteTag("foo", "bar", NilAuth{})

	assert.NoError(t, err)
}

func TestRepositoryDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repositories/foo/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
	})

	err := client.Repository.Delete("foo", NilAuth{})

	assert.NoError(t, err)
}
