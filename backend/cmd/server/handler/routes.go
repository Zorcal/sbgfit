package handler

import (
	"fmt"
	"net/http"

	"github.com/zorcal/sbgfit/backend/pkg/httprouter"
)

func routes(r *httprouter.Router) {
	r.Handle("GET /{$}", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "Hello world!")
		return nil
	})
}
