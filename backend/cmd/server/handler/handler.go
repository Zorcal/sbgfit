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

func respond[T any](w http.ResponseWriter, statusCode int, data T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if statusCode == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.WriteHeader(statusCode)

	envelope := struct {
		Data T `json:"data"`
	}{
		Data: data,
	}
	if err := json.NewEncoder(w).Encode(envelope); err != nil {
		return fmt.Errorf("respond: encode json: %w", err)
	}

	return nil
}
