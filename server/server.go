package server // import "github.com/CenturyLinkLabs/imagelayers/server"

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CenturyLinkLabs/imagelayers/api"
	"github.com/CenturyLinkLabs/imagelayers/graph"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

type statusLoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusLoggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (s *layerServer) createRouter() serverRouter {
	registry := api.NewRemoteRegistry()
	router := serverRouter{mux.NewRouter()}

	router.addCorsRoutes("/registry", registry.Routes())
	router.addRoutes("/graph", graph.Routes())

	return router
}

type serverRouter struct {
	*mux.Router
}

func (sr *serverRouter) addCorsRoutes(context string, routeMap map[string]map[string]http.HandlerFunc) {
	sr.generateRoutes(context, routeMap, func(path string, method string, wrap http.HandlerFunc) {
		sr.Path(path).Methods(method).Handler(cors.Default().Handler(wrap))
		sr.Path(path).Methods("OPTIONS").Handler(cors.Default().Handler(wrap))
	})
}

func (sr *serverRouter) addRoutes(context string, routeMap map[string]map[string]http.HandlerFunc) {
	sr.generateRoutes(context, routeMap, func(path string, method string, wrap http.HandlerFunc) {
		sr.Path(path).Methods(method).Handler(wrap)
	})
}

func (sr *serverRouter) generateRoutes(context string, routeMap map[string]map[string]http.HandlerFunc, build func(string, string, http.HandlerFunc)) {

	for method, routes := range routeMap {
		for route, fct := range routes {

			localMethod := method
			localRoute := route
			localFct := fct
			if localMethod == "STATIC" {
				sr.generateStaticRoute(context, localRoute, localFct)
			} else {
				wrap := func(w http.ResponseWriter, r *http.Request) {
					ww := &statusLoggingResponseWriter{w, 200}

					log.Printf("Started %s %s", r.Method, r.RequestURI)

					localFct(ww, r)

					log.Printf("Completed %d", ww.statusCode)
				}

				build(context+localRoute, localMethod, wrap)
			}
		}
	}
}

func (sr *serverRouter) generateStaticRoute(context string, route string, staticHandler http.HandlerFunc) {
	path := context + route
	log.Printf("static: %s", path)
	sr.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir("."+path))))
}
