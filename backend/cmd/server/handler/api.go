package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/zorcal/sbgfit/backend/cmd/server/handler/openapi"
	"github.com/zorcal/sbgfit/backend/core/mdl"
)

//go:generate moq -rm -fmt goimports -out exercise_service_moq_test.go . ExerciseService:MockeExerciseServiced

type ExerciseService interface {
	Exercises(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error)
}

type api struct {
	log         *slog.Logger
	exerciseSvc ExerciseService
}

func (a *api) NewError(ctx context.Context, err error) *openapi.ErrorResponseStatusCode {
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
