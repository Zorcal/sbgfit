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
		"barbell",
		"dumbbell",
		"kettlebell",
		"cable",
		"machine",
		"bodyweight",
		"bands",
		"medicine_ball",
		"plates",
		"pull_up_bar",
		"box",
		"rope",
		"sled",
		"trx",
		"other",
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
		"quads",
		"hamstrings",
		"glutes",
		"calves",
		"abs",
		"traps",
		"rhomboids",
		"lats",
		"forearms",
		"cardio",
		"full_body",
	}
	return slices.Contains(validMuscles, s)
}

// IsValidExerciseTag checks if the given string is a valid exercise tag.
func IsValidExerciseTag(s string) bool {
	validTags := []string{
		"crossfit",
		"hyrox",
		"bodybuilding",
		"strength",
		"powerlifting",
		"olympic_lifting",
		"cardio",
		"hiit",
		"yoga",
		"pilates",
		"stretching",
		"rehab",
		"beginner",
		"intermediate",
		"advanced",
		"compound",
		"isolation",
		"unilateral",
		"bilateral",
		"explosive",
		"isometric",
	}
	return slices.Contains(validTags, s)
}
