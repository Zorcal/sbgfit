package api

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/zorcal/sbgfit/backend/api/internal/conv"
	"github.com/zorcal/sbgfit/backend/api/internal/openapi"
	"github.com/zorcal/sbgfit/backend/internal/core/mdl"
	"github.com/zorcal/sbgfit/backend/internal/telemetry"
	"github.com/zorcal/sbgfit/backend/pkg/slicesx"
)

//go:generate moq -rm -fmt goimports -pkg api_test -out exercise_service_moq_test.go . ExerciseService:MockedExerciseServiced

type ExerciseService interface {
	Exercises(ctx context.Context, fltr mdl.ExerciseFilter, pageSize, pageNumber int) (exs []mdl.Exercise, totalCount int, err error)
}

func (a *api) GetExercises(ctx context.Context, params openapi.GetExercisesParams) (openapi.GetExercisesRes, error) {
	ctx, span := telemetry.StartSpan(ctx, "api.api.GetExercises")
	defer span.End()

	span.SetAttributes(exercisesParamsSpanAttributes(params)...)

	fltr := conv.ExerciseFilterFromAPI(params)

	pageSize := 20
	if ps, ok := params.PageSize.Get(); ok {
		pageSize = ps
	}

	pageNumber := 1
	if pn, ok := params.PageNumber.Get(); ok {
		pageNumber = pn
	}

	exs, totalCount, err := a.exerciseSvc.Exercises(ctx, fltr, pageSize, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("get exercises: %w", err)
	}

	data := slicesx.Map(exs, conv.ExerciseToAPI)

	return &openapi.ExerciseResponse{
		Data:  data,
		Total: totalCount,
	}, nil
}

func exercisesParamsSpanAttributes(params openapi.GetExercisesParams) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.Int("exercise_params.page_size", params.PageSize.Value),
		attribute.Int("exercise_params.page_number", params.PageNumber.Value),
	}

	if name, ok := params.Name.Get(); ok {
		attrs = append(attrs, attribute.String("exercise_params.name", name))
	}

	if category, ok := params.Category.Get(); ok {
		attrs = append(attrs, attribute.String("exercise_params.category", string(category)))
	}

	if len(params.EquipmentTypes) > 0 {
		attrs = append(attrs, attribute.StringSlice("exercise_params.equipment_types", slicesx.ToStrings(params.EquipmentTypes)))
	}

	if len(params.PrimaryMuscles) > 0 {
		attrs = append(attrs, attribute.StringSlice("exercise_params.primary_muscles", slicesx.ToStrings(params.PrimaryMuscles)))
	}

	if len(params.Tags) > 0 {
		attrs = append(attrs, attribute.StringSlice("exercise_params.tags", slicesx.ToStrings(params.Tags)))
	}

	return attrs
}
