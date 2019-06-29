package lemur

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sjsafranek/ligneous"
)

const DEFAULT_HTTP_PORT = 8080

func NewServer(port int, log ligneous.Log) (*HttpServer, error) {
	return &HttpServer{Router: mux.NewRouter().StrictSlash(true), Port: port, Log: log}, nil
}

type HttpServer struct {
	Port   int
	Router *mux.Router
	Log    ligneous.Log
}

func (self *HttpServer) AttachHandlers(routes []ApiRoute) {
	for _, route := range routes {
		self.AttachHandler(route)
	}
}

func (self *HttpServer) AttachHandler(route ApiRoute) {
	var handler http.Handler
	self.Log.Infof("Attaching HTTP handler for route: %v %v", route.Methods, route.Pattern)
	handler = route.HandlerFunc
	self.Router.
		Methods(route.Methods...).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
}

func (self *HttpServer) Listen() {
	self.Log.Info(fmt.Sprintf("Magic happens on port %v", self.Port))

	bind := fmt.Sprintf(":%v", self.Port)

	self.Router.Use(LoggingMiddleWare(self.Log), SetHeadersMiddleWare, CORSMiddleWare)

	err := http.ListenAndServe(bind, self.Router)
	if err != nil {
		panic(err)
	}
}
