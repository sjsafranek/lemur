package lemur

import (
	"net/http"
)

type ApiRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
