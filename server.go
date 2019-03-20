package lemur

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const DEFAULT_HTTP_PORT = 8080

type HttpServer struct {
	Port   int
	Router *mux.Router
}

func (self HttpServer) Start() {
	// Start server
	logger.Info(fmt.Sprintf("Magic happens on port %v", self.Port))

	bind := fmt.Sprintf(":%v", self.Port)

	self.Router.Use(LoggingMiddleWare, SetHeadersMiddleWare, CORSMiddleWare)

	err := http.ListenAndServe(bind, self.Router)
	if err != nil {
		panic(err)
	}

}
