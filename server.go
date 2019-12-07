package lemur

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sjsafranek/logger"
)

const DEFAULT_HTTP_PORT = 8080

func NewServer() (*HttpServer, error) {
	return &HttpServer{Router: mux.NewRouter().StrictSlash(true)}, nil
}

type HttpServer struct {
	Router *mux.Router
}

func (self *HttpServer) AttachHandlerFuncs(routes []ApiRoute) {
	for _, route := range routes {
		self.AttachHandlerFunc(route)
	}
}

func (self *HttpServer) AttachHandlerFunc(route ApiRoute) {
	logger.Infof("Attaching HTTP handler for route: %v %v", route.Methods, route.Pattern)
	self.Router.
		Methods(route.Methods...).
		Path(route.Pattern).
		Name(route.Name).
		Handler(route.HandlerFunc)
}

func (self *HttpServer) AttachFileServer(path, directory string) {
	fsvr := http.FileServer(http.Dir(directory))
	self.Router.
		PathPrefix(path).
		Handler(http.StripPrefix(path, fsvr))
}

func (self *HttpServer) AttachHandler(path string, handler http.Handler) {
	self.Router.PathPrefix(path).Handler(handler)
}

func (self *HttpServer) ListenAndServe(port int) {
	logger.Info(fmt.Sprintf("Magic happens on port %v", port))

	bind := fmt.Sprintf(":%v", port)

	self.Router.Use(LoggingMiddleWare(), SetHeadersMiddleWare, CORSMiddleWare)

	err := http.ListenAndServe(bind, self.Router)
	if err != nil {
		panic(err)
	}
}
