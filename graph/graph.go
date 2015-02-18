package graph // import "github.com/CenturyLinkLabs/imagelayers/graph"

import (
	"net/http"
	"log"
	"text/template"
)

type Test struct {
	Name string `json:"name"`
}

func Routes() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"GET": {
			"/test":   indexHandler,
		},
		"STATIC": {
			"/assets": assetHandler,
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

func assetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("static url: %s", r.URL)
}
