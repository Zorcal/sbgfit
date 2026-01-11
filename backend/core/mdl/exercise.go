package mdl

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

// ExerciseFilter represents search and filtering criteria for finding
// exercises from the exercise library that match specific training
// requirements. Used to help users discover exercises they might want to clone
// for their personal use.
type ExerciseFilter struct {
	Name           *string
	Category       *string
	EquipmentTypes []string
	PrimaryMuscles []string
	Tags           []string
}

// Exercise represents a standardized training movement from the exercise
// library provided by the application. These exercises are static, predefined
// movements that serve as templates for users. Users can clone exercises from
// this library to create their own customizable "user exercises" that they can
// modify to fit their specific needs and preferences.
type Exercise struct {
	ID             uuid.UUID
	Name           string
	Category       string
	Description    *string
	Instructions   []string
	EquipmentTypes []string
	PrimaryMuscles []string
	Tags           []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
