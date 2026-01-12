package api

import (
	"context"
	"fmt"

	"github.com/zorcal/sbgfit/backend/api/adapt"
	"github.com/zorcal/sbgfit/backend/api/openapi"
	"github.com/zorcal/sbgfit/backend/pkg/slicesx"
)

func (a *api) GetExercises(ctx context.Context, params openapi.GetExercisesParams) (openapi.GetExercisesRes, error) {
	fltr := adapt.ExerciseFilterFromAPI(params)

	exs, err := a.exerciseSvc.Exercises(ctx, fltr)
	if err != nil {
		return nil, fmt.Errorf("get exercises: %w", err)
	}

	data := slicesx.Map(exs, adapt.ExerciseToAPI)

	return &openapi.ExerciseResponse{
		Data: data,
	}, nil
}
