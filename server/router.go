package server // import "github.com/CenturyLinkLabs/imagelayers/server"

import (
	"expvar"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	expRequests = expvar.NewInt("requests")
)

type statusLoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusLoggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

type Router struct {
	*mux.Router
}

type RouteMap map[string]map[string]http.HandlerFunc

func (sr *Router) AddCorsRoutes(context string, routeMap RouteMap) {
	sr.generateRoutes(context, routeMap, func(path string, method string, wrap http.HandlerFunc) {
		corsHndlr := cors.New(cors.Options{
			AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
			AllowCredentials: true,
		})

		sr.Path(path).Methods("OPTIONS").Handler(corsHndlr.Handler(wrap))
		sr.Path(path).Methods(method).Handler(corsHndlr.Handler(wrap))
	})
}

func (sr *Router) AddRoutes(context string, routeMap RouteMap) {
	sr.generateRoutes(context, routeMap, func(path string, method string, wrap http.HandlerFunc) {
		sr.Path(path).Methods(method).Handler(wrap)
	})
}

func (sr *Router) generateRoutes(context string, routeMap RouteMap, build func(string, string, http.HandlerFunc)) {

	for method, routes := range routeMap {
		for route, fct := range routes {

			localMethod := method
			localRoute := route
			localFct := fct

			wrap := func(w http.ResponseWriter, r *http.Request) {
				ww := &statusLoggingResponseWriter{w, 200}

				expRequests.Add(1)
				log.Printf("Started %s %s", r.Method, r.RequestURI)

				localFct(ww, r)

				log.Printf("Completed %d", ww.statusCode)
			}

			build(context+localRoute, localMethod, wrap)

		}
	}
}
