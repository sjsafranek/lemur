package lemur

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/schollz/httpfileserver"
	"github.com/sjsafranek/lemur/middleware"
	"github.com/sjsafranek/logger"
)

func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func Body(r *http.Request, clbk func([]byte) error) error {
	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		return err
	}
	defer r.Body.Close()
	return clbk(body)
}

type HttpRouter struct {
	*mux.Router
}

func (self *HttpRouter) AttachFileServer(route, directory string) {
	self.PathPrefix("/static/").Handler(httpfileserver.New("/static/", directory))
}

func (self *HttpRouter) AttachHandlerFunc(pattern string, handler http.HandlerFunc, methods []string) {
	logger.Info("Attaching HTTP handler for route: ", methods, " ", pattern)
	self.Methods(methods...).Path(pattern).Handler(handler)
}

func New() HttpRouter {
	// create http router
	router := mux.NewRouter().StrictSlash(true)

	// attach middleware
	router.Use(middleware.LoggingMiddleWare, middleware.SetHeadersMiddleWare, middleware.CORSMiddleWare)
	handlers.CompressHandler(router)

	return HttpRouter{router}
}
