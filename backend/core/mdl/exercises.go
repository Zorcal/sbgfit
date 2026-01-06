package mdl

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

// ExerciseFilter represents search and filtering criteria for finding exercises
// that match specific training requirements and user preferences.
type ExerciseFilter struct {
	Name           *string
	Category       *string
	EquipmentTypes []string
	PrimaryMuscles []string
	Tags           []string
	CreatedByUser  *bool
}

// Exercise represents a standardized training movement that forms the building
// blocks of workouts. Exercises are the atomic units of training that users
// select and combine to create personalized workout routines.
type Exercise struct {
	ID              uuid.UUID
	Name            string
	Category        string
	Description     *string
	Instructions    []string
	EquipmentTypes  []string
	PrimaryMuscles  []string
	Tags            []string
	CreatedByUserID uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// IsValidEquipmentType checks if the given string is a valid equipment type.
func IsValidEquipmentType(s string) bool {
	validTypes := []string{
		"bodyweight",
		"kettlebell",
		"rowing-machine",
		"ski-erg",
		"medicine-ball",
		"dumbbells",
		"barbell",
		"sled",
		"box",
		"jump-rope",
		"assault-bike",
	}
	return slices.Contains(validTypes, s)
}

// IsValidPrimaryMuscle checks if the given string is a valid primary muscle group.
func IsValidPrimaryMuscle(s string) bool {
	validMuscles := []string{
		"chest",
		"back",
		"shoulders",
		"biceps",
		"triceps",
		"forearms",
		"core",
		"abs",
		"obliques",
		"glutes",
		"quads",
		"hamstrings",
		"calves",
		"legs",
		"full-body",
		"grip",
	}
	return slices.Contains(validMuscles, s)
}

// IsValidExerciseTag checks if the given string is a valid exercise tag.
func IsValidExerciseTag(s string) bool {
	validTags := []string{
		"crossfit",
		"hyrox",
		"beginner-friendly",
		"advanced",
		"conditioning",
		"strength-endurance",
		"power",
		"core",
		"functional",
		"competition",
		"plyometric",
	}
	return slices.Contains(validTags, s)
}

// IsValidCategory checks if the given string is a valid exercise category.
func IsValidCategory(s string) bool {
	validCategories := []string{
		"cardio",
		"strength",
		"plyometric",
	}
	return slices.Contains(validCategories, s)
}
