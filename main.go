package main // import "github.com/microscaling/imagelayers"

import (
	"expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/microscaling/imagelayers/api"
	"github.com/microscaling/imagelayers/server"
	"gopkg.in/tylerb/graceful.v1"
)

type layerServer struct {
}

func NewServer() *layerServer {
	return new(layerServer)
}

func (s *layerServer) Start(port int) {
	router := s.createRouter()

	log.Printf("Server starting on port %d", port)
	portString := fmt.Sprintf(":%d", port)
	graceful.Run(portString, 10*time.Second, router)
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
	s.Start(*port)
}
