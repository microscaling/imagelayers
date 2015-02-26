package graph // import "github.com/CenturyLinkLabs/imagelayers/graph"

import (
	"net/http"
	"text/template"
	"github.com/CenturyLinkLabs/imagelayers/server"
)

type Test struct {
	Name string `json:"name"`
}

func Routes(context string, router *server.Router) {
	routes :=  map[string]map[string]http.HandlerFunc{
		"GET": {
			"/":   indexHandler,
		},
	}

	router.AddRoutes(context, routes)

	// Add Static Route
	path := context
	router.PathPrefix(path).Handler(http.FileServer(http.Dir("./graph")))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html")
	t, err := t.ParseFiles("graph/index.html")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, "")
}
