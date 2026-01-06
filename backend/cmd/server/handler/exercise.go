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
		fltr, err := buildExerciseFilter(r)
		if err != nil {
			return &httpError{
				StatusCode:      http.StatusBadRequest,
				ExternalMessage: "Invalid filter query params",
				InternalErr:     err,
			}
		}

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

func buildExerciseFilter(r *http.Request) (mdl.ExerciseFilter, error) {
	var filter mdl.ExerciseFilter

	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = ptr.To(name)
	}

	if category := r.URL.Query().Get("category"); category != "" {
		if !mdl.IsValidCategory(category) {
			return mdl.ExerciseFilter{}, fmt.Errorf("invalid category: %q", category)
		}
		filter.Category = ptr.To(category)
	}

	if equipmentTypes := r.URL.Query().Get("equipmentTypes"); equipmentTypes != "" {
		types := strings.Split(equipmentTypes, ",")
		for _, t := range types {
			if t != "" && !mdl.IsValidEquipmentType(t) {
				return mdl.ExerciseFilter{}, fmt.Errorf("invalid equipment type: %q", t)
			}
		}
		filter.EquipmentTypes = types
	}

	if primaryMuscles := r.URL.Query().Get("primaryMuscles"); primaryMuscles != "" {
		muscles := strings.Split(primaryMuscles, ",")
		for _, m := range muscles {
			if m != "" && !mdl.IsValidPrimaryMuscle(m) {
				return mdl.ExerciseFilter{}, fmt.Errorf("invalid primary muscle: %q", m)
			}
		}
		filter.PrimaryMuscles = muscles
	}

	if tags := r.URL.Query().Get("tags"); tags != "" {
		tagList := strings.Split(tags, ",")
		for _, tag := range tagList {
			if tag != "" && !mdl.IsValidExerciseTag(tag) {
				return mdl.ExerciseFilter{}, fmt.Errorf("invalid exercise tag: %q", tag)
			}
		}
		filter.Tags = tagList
	}

	if createdByUser := r.URL.Query().Get("createdByUser"); createdByUser != "" {
		if val, err := strconv.ParseBool(createdByUser); err == nil {
			filter.CreatedByUser = ptr.To(val)
		} else {
			return mdl.ExerciseFilter{}, fmt.Errorf("invalid createdByUser value: %q", createdByUser)
		}
	}

	return filter, nil
}
