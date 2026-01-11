// Package handler provides the handler for the HTTP server.
package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/middleware"

	"github.com/zorcal/sbgfit/backend/cmd/server/handler/openapi"
)

type Config struct {
	Log             *slog.Logger
	ExerciseService ExerciseService
}

func New(cfg Config) (http.Handler, error) {
	srv, err := openapi.NewServer(
		&api{
			log:         cfg.Log,
			exerciseSvc: cfg.ExerciseService,
		},
		openapi.WithMiddleware(middleware.ChainMiddlewares(
			panicRecoveryMiddleware(cfg.Log),
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("create openapi server: %w", err)
	}

	h := wrapHTTPMiddleware(srv,
		httpTraceMiddleware(),
		httpLoggingMiddleware(cfg.Log),
	)

	return h, nil
}
