package lemur

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/schollz/httpfileserver"
	"github.com/sjsafranek/lemur/middleware"
	"github.com/sjsafranek/logger"
)

var (
	cacheSince = time.Now().Format(http.TimeFormat)
	cacheUntil = time.Now().AddDate(0, 0, 1).Format(http.TimeFormat)
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

func CacheControlWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=86400, must-revalidate, public") // 30 days
		w.Header().Set("Last-Modified", cacheSince)
		w.Header().Set("Expires", cacheUntil)
		h.ServeHTTP(w, r)
	})
}

func (self *HttpRouter) AttachFileServer(route, directory string) {
	fs := CacheControlWrapper(httpfileserver.New("/static/", directory))
	self.PathPrefix("/static/").Handler(fs)
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
