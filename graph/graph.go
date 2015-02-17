package graph // import "github.com/CenturyLinkLabs/imagelayers/graph"

import (
	"encoding/json"
	"net/http"
)

type Test struct {
	Name string `json:"name"`
}

func Routes() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"GET": {
			"/": testGraph,
		},
	}
}

func testGraph(w http.ResponseWriter, r *http.Request) {
	test := new(Test)
	test.Name = "GRAPH"

	json.NewEncoder(w).Encode(test)
}
