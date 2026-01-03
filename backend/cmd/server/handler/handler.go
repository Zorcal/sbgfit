// Package handler provides the handler for the HTTP server.
package handler

import (
	"log/slog"
	"net/http"

	"github.com/zorcal/sbgfit/backend/pkg/httprouter"
)

type Config struct {
	Log *slog.Logger
}

func New(cfg Config) http.Handler {
	r := httprouter.New(
		traceMiddleware(),
		loggingMiddleware(cfg.Log),
		panicRecovery(cfg.Log),
	)
	routes(r)
	return r
}
