// Package handler provides the handler for the HTTP server.
package handler

import (
	"fmt"
	"net/http"
)

func New() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world!")
	})
}
