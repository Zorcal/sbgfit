package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/middleware"
	"github.com/zorcal/sbgfit/backend/api/openapi"
)

type Config struct {
	Log             *slog.Logger
	ExerciseService ExerciseService
}

func NewHandler(cfg Config) (http.Handler, error) {
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
