package main // import "github.com/CenturyLinkLabs/imagelayers"

import (
	"expvar"
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

	// Handler for ExpVar request
	router.Path("/debug/vars").Methods("GET").HandlerFunc(expvarHandler)

	return router
}

func expvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
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
