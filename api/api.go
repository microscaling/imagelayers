package api

import (
	"encoding/json"
	"net/http"
)

type Test struct {
	Name string `json:"name"`
}

type ImageManager interface {
	Test() (Test, error)
}

func Routes() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc {
		"GET": {
			"/test": testApi,
		},
	}
}

func testApi(w http.ResponseWriter, r *http.Request) {
	test := new(Test)
	test.Name = "Me"

	json.NewEncoder(w).Encode(test)
}
