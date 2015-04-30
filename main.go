package main // import "github.com/CenturyLinkLabs/imagelayers"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CenturyLinkLabs/imagelayers/api"
	"github.com/CenturyLinkLabs/imagelayers/server"
	"github.com/gorilla/mux"
)

type layerServer struct {
}

func NewServer() *layerServer {
	return new(layerServer)
}

func (s *layerServer) Start(port int) error {
	router := s.createRouter()

	log.Printf("Server starting on port %d", port)
	portString := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(portString, router)
}

func (s *layerServer) createRouter() server.Router {
	registry := api.NewRemoteRegistry()
	router := server.Router{mux.NewRouter()}

	registry.Routes("/registry", &router)

	return router
}

func main() {
	port := flag.Int("p", 8888, "port on which the server will run")
	flag.Parse()

	s := NewServer()
	if err := s.Start(*port); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
