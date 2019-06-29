package lemur

import (
	"net/http"
)

type ApiRoute struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
