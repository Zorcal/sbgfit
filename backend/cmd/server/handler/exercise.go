package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zorcal/sbgfit/backend/core/mdl"
	"github.com/zorcal/sbgfit/backend/pkg/httprouter"
	"github.com/zorcal/sbgfit/backend/pkg/ptr"
)

//go:generate moq -rm -fmt goimports -out exercise_service_moq_test.go . ExerciseService:MockeExerciseServiced

type ExerciseService interface {
	Exercises(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error)
}

func getExercisesHandler(svc ExerciseService) httprouter.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fltr := buildExerciseFilter(r)
		data, err := svc.Exercises(r.Context(), fltr)
		if err != nil {
			return fmt.Errorf("get exercise: %w", err)
		}

		if err := respond(w, http.StatusOK, data); err != nil {
			return fmt.Errorf("respond: %w", err)
		}

		return nil
	}
}

func buildExerciseFilter(r *http.Request) mdl.ExerciseFilter {
	var filter mdl.ExerciseFilter

	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = ptr.To(name)
	}

	if category := r.URL.Query().Get("category"); category != "" {
		filter.Category = ptr.To(category)
	}

	if equipmentTypes := r.URL.Query().Get("equipmentTypes"); equipmentTypes != "" {
		filter.EquipmentTypes = strings.Split(equipmentTypes, ",")
	}

	if primaryMuscles := r.URL.Query().Get("primaryMuscles"); primaryMuscles != "" {
		filter.PrimaryMuscles = strings.Split(primaryMuscles, ",")
	}

	if tags := r.URL.Query().Get("tags"); tags != "" {
		filter.Tags = strings.Split(tags, ",")
	}

	if createdByUser := r.URL.Query().Get("createdByUser"); createdByUser != "" {
		if val, err := strconv.ParseBool(createdByUser); err == nil {
			filter.CreatedByUser = ptr.To(val)
		}
	}

	return filter
}
