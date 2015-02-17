package graph // import "github.com/CenturyLinkLabs/imagelayers/graph"

import (
	"net/http"
	"text/template"
)

type Test struct {
	Name string `json:"name"`
}

func Routes() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"GET": {
			"/": indexHandler,
		},
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html")
	t, err := t.ParseFiles("graph/tmpl/index.html")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, "")

}
