package api

import (
	"context"
	"fmt"

	"github.com/zorcal/sbgfit/backend/api/internal/conv"
	"github.com/zorcal/sbgfit/backend/api/internal/openapi"
	"github.com/zorcal/sbgfit/backend/internal/core/mdl"
	"github.com/zorcal/sbgfit/backend/pkg/slicesx"
)

//go:generate moq -rm -fmt goimports -pkg api_test -out exercise_service_moq_test.go . ExerciseService:MockedExerciseServiced

type ExerciseService interface {
	Exercises(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error)
}

func (a *api) GetExercises(ctx context.Context, params openapi.GetExercisesParams) (openapi.GetExercisesRes, error) {
	fltr := conv.ExerciseFilterFromAPI(params)

	exs, err := a.exerciseSvc.Exercises(ctx, fltr)
	if err != nil {
		return nil, fmt.Errorf("get exercises: %w", err)
	}

	data := slicesx.Map(exs, conv.ExerciseToAPI)

	return &openapi.ExerciseResponse{
		Data: data,
	}, nil
}
