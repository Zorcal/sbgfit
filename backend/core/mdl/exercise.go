package mdl

import (
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
