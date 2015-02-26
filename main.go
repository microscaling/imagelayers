package main // import "github.com/CenturyLinkLabs/imagelayers"

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/CenturyLinkLabs/imagelayers/server"
	"github.com/CenturyLinkLabs/imagelayers/api"
	"github.com/CenturyLinkLabs/imagelayers/graph"
)

type layerServer struct {
}

func NewServer() *layerServer {
	return new(layerServer)
}

func (s *layerServer) Start(port int) {
	router := s.createRouter()

	log.Printf("Server running on port %d", port)
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, router)
}

func (s *layerServer) createRouter() server.Router {
	registry := api.NewRemoteRegistry()
	router := server.Router{mux.NewRouter()}

	registry.Routes("/registry", &router)
	graph.Routes("/", &router)

	return router
}

func main() {
	port := flag.Int("p", 9000, "port on which the server will run")
	flag.Parse()

	s := NewServer()
	s.Start(*port)
}
