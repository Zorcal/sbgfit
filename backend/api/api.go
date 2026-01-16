// Package api provides the API handler for the HTTP server.
package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/codes"

	"github.com/zorcal/sbgfit/backend/api/internal/openapi"
	"github.com/zorcal/sbgfit/backend/internal/telemetry"
)

type api struct {
	log         *slog.Logger
	exerciseSvc ExerciseService
}

func (a *api) NewError(ctx context.Context, err error) *openapi.ErrorResponseStatusCode {
	span := telemetry.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	if httpErr := new(httpError); errors.As(err, &httpErr) {
		a.log.Log(ctx, logLevel(httpErr.StatusCode), "Request error", "error", httpErr)
		return &openapi.ErrorResponseStatusCode{
			StatusCode: httpErr.StatusCode,
			Response: openapi.ErrorResponse{
				Error: httpErr.ExternalMessage,
			},
		}
	}

	a.log.ErrorContext(ctx, "Request error", "error", err)
	return &openapi.ErrorResponseStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: openapi.ErrorResponse{
			Error: http.StatusText(http.StatusInternalServerError),
		},
	}
}

func logLevel(statusCode int) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError
	case statusCode >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
