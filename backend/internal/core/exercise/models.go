package exercise

import (
	"time"

	"github.com/google/uuid"

	"github.com/zorcal/sbgfit/backend/internal/core/mdl"
)

type dbExercisesResult struct {
	dbExercise

	TotalCount int `db:"total_count"`
}

type dbExercise struct {
	ExternalID     uuid.UUID `db:"external_id"`
	Name           string    `db:"name"`
	CategoryCode   string    `db:"category_code"`
	Description    *string   `db:"description"`
	Instructions   []string  `db:"instructions"`
	EquipmentTypes []string  `db:"equipment_types"`
	PrimaryMuscles []string  `db:"primary_muscles"`
	Tags           []string  `db:"tags"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func dbExerciseToModel(db dbExercise) mdl.Exercise {
	return mdl.Exercise{
		ID:             db.ExternalID,
		Name:           db.Name,
		Category:       db.CategoryCode,
		Description:    db.Description,
		Instructions:   db.Instructions,
		EquipmentTypes: db.EquipmentTypes,
		PrimaryMuscles: db.PrimaryMuscles,
		Tags:           db.Tags,
		CreatedAt:      db.CreatedAt,
		UpdatedAt:      db.UpdatedAt,
	}
}
