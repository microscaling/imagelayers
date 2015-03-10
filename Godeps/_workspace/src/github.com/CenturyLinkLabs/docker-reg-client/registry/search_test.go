package registry

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchQuery(t *testing.T) {
	setup()
	defer teardown()

	expected := &SearchResults{
		NumPages:   1,
		NumResults: 1,
		Results:    []SearchResult{{Name: "foo"}},
		PageSize:   1,
		Query:      "foo",
		Page:       1,
	}

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "foo", r.FormValue("q"))
		assert.Equal(t, "2", r.FormValue("page"))
		assert.Equal(t, "3", r.FormValue("n"))
		json.NewEncoder(w).Encode(expected)
	})

	results, err := client.Search.Query("foo", 2, 3)

	assert.NoError(t, err)
	assert.Equal(t, expected, results)
}
