package api

import (
	"context"
	"fmt"

	"github.com/zorcal/sbgfit/backend/api/internal/conv"
	"github.com/zorcal/sbgfit/backend/api/internal/openapi"
	"github.com/zorcal/sbgfit/backend/pkg/slicesx"
)

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
