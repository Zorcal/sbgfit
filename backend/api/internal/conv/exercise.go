package conv

import (
	"github.com/zorcal/sbgfit/backend/api/internal/openapi"
	"github.com/zorcal/sbgfit/backend/internal/core/mdl"
	"github.com/zorcal/sbgfit/backend/pkg/ptr"
	"github.com/zorcal/sbgfit/backend/pkg/slicesx"
)

func ExerciseToAPI(ex mdl.Exercise) openapi.Exercise {
	var description openapi.OptNilString
	if ex.Description != nil {
		description.SetTo(*ex.Description)
	} else {
		description.SetToNull()
	}

	return openapi.Exercise{
		ID:             ex.ID,
		Name:           ex.Name,
		Category:       openapi.ExerciseCategory(ex.Category),
		Description:    description,
		Instructions:   ex.Instructions,
		EquipmentTypes: slicesx.Map(ex.EquipmentTypes, func(s string) openapi.EquipmentType { return openapi.EquipmentType(s) }),
		PrimaryMuscles: slicesx.Map(ex.PrimaryMuscles, func(s string) openapi.PrimaryMuscle { return openapi.PrimaryMuscle(s) }),
		Tags:           slicesx.Map(ex.Tags, func(s string) openapi.ExerciseTag { return openapi.ExerciseTag(s) }),
		CreatedAt:      ex.CreatedAt,
		UpdatedAt:      ex.UpdatedAt,
	}
}

func ExerciseFromAPI(ex openapi.Exercise) mdl.Exercise {
	var description *string
	if desc, ok := ex.Description.Get(); ok {
		description = &desc
	}

	return mdl.Exercise{
		ID:             ex.ID,
		Name:           ex.Name,
		Category:       string(ex.Category),
		Description:    description,
		Instructions:   ex.Instructions,
		EquipmentTypes: slicesx.Map(ex.EquipmentTypes, func(e openapi.EquipmentType) string { return string(e) }),
		PrimaryMuscles: slicesx.Map(ex.PrimaryMuscles, func(m openapi.PrimaryMuscle) string { return string(m) }),
		Tags:           slicesx.Map(ex.Tags, func(t openapi.ExerciseTag) string { return string(t) }),
		CreatedAt:      ex.CreatedAt,
		UpdatedAt:      ex.UpdatedAt,
	}
}

func ExerciseFilterFromAPI(params openapi.GetExercisesParams) mdl.ExerciseFilter {
	var filter mdl.ExerciseFilter

	if name, ok := params.Name.Get(); ok {
		filter.Name = ptr.To(name)
	}

	if category, ok := params.Category.Get(); ok {
		filter.Category = ptr.To(string(category))
	}

	if len(params.EquipmentTypes) > 0 {
		filter.EquipmentTypes = slicesx.Map(params.EquipmentTypes, func(e openapi.EquipmentType) string { return string(e) })
	}

	if len(params.PrimaryMuscles) > 0 {
		filter.PrimaryMuscles = slicesx.Map(params.PrimaryMuscles, func(m openapi.PrimaryMuscle) string { return string(m) })
	}

	if len(params.Tags) > 0 {
		filter.Tags = slicesx.Map(params.Tags, func(t openapi.ExerciseTag) string { return string(t) })
	}

	return filter
}
