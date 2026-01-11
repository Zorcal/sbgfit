// Package handler provides the handler for the HTTP server.
package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/zorcal/sbgfit/backend/pkg/httprouter"
)

type Config struct {
	Log             *slog.Logger
	ExerciseService ExerciseService
}

func New(cfg Config) http.Handler {
	r := httprouter.New(
		traceMiddleware(),
		loggingMiddleware(cfg.Log),
		errorMiddleware(cfg.Log),
		panicRecovery(cfg.Log),
	)
	routes(r, cfg)
	return r
}

type response[T any] struct {
	Data T `json:"data"`
}

func respond[T any](w http.ResponseWriter, statusCode int, data T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if statusCode == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.WriteHeader(statusCode)

	resp := response[T]{
		Data: data,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
