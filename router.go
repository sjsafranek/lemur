package lemur

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router for http api calls
func Router(routes []ApiRoute) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	AttachHandlers(router, routes)
	return router
}

func AttachHandlers(router *mux.Router, routes []ApiRoute) {
	for _, route := range routes {
		AttachHandler(router, route)
	}
}

func AttachHandler(router *mux.Router, route ApiRoute) {
	var handler http.Handler
	logger.Info("Attaching HTTP handler for route: ", route.Method, " ", route.Pattern)
	handler = route.HandlerFunc
	router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
}
