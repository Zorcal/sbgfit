// Package handler provides the handler for the HTTP server.
package handler

import (
	"net/http"

	"github.com/zorcal/sbgfit/backend/pkg/httprouter"
)

func New() http.Handler {
	r := httprouter.New()
	routes(r)
	return r
}
