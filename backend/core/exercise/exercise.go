// Package exercise provides the application service for managing exercises.
package exercise

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zorcal/sbgfit/backend/core/mdl"
	"github.com/zorcal/sbgfit/backend/pkg/ptr"
)

// Service is the exercise application service.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a new exercise service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		pool: pool,
	}
}

// Exercises retrieves a list of exercises based on the provided filter criteria.
func (s *Service) Exercises(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error) {
	now := time.Now()

	exercises := []mdl.Exercise{
		{
			ID:              uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			Name:            "Deadlift",
			Category:        "strength",
			Description:     ptr.To("A compound hip-hinge movement that targets the posterior chain"),
			Instructions:    []string{"Stand with feet hip-width apart", "Grip bar with hands outside legs", "Lift by driving hips forward"},
			EquipmentTypes:  []string{"barbell", "plates"},
			PrimaryMuscles:  []string{"hamstrings", "glutes", "back"},
			Tags:            []string{"compound", "powerlifting", "crossfit"},
			CreatedByUserID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			CreatedAt:       now.Add(-24 * time.Hour),
			UpdatedAt:       now.Add(-24 * time.Hour),
		},
		{
			ID:              uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Name:            "Burpees",
			Category:        "cardio",
			Description:     ptr.To("Full-body exercise combining a squat thrust with a jump"),
			Instructions:    []string{"Start standing", "Drop into squat position", "Jump back to plank", "Do push-up", "Jump feet forward", "Jump up with arms overhead"},
			EquipmentTypes:  []string{"bodyweight"},
			PrimaryMuscles:  []string{"full_body", "cardio"},
			Tags:            []string{"hiit", "crossfit", "hyrox", "bodyweight"},
			CreatedByUserID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			CreatedAt:       now.Add(-48 * time.Hour),
			UpdatedAt:       now.Add(-48 * time.Hour),
		},
		{
			ID:              uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			Name:            "Bicep Curls",
			Category:        "strength",
			Description:     ptr.To("Isolated arm exercise targeting the bicep muscles"),
			Instructions:    []string{"Hold dumbbells at sides", "Keep elbows at sides", "Curl weights up to shoulders", "Lower with control"},
			EquipmentTypes:  []string{"dumbbell"},
			PrimaryMuscles:  []string{"biceps"},
			Tags:            []string{"isolation", "bodybuilding", "beginner"},
			CreatedByUserID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			CreatedAt:       now.Add(-72 * time.Hour),
			UpdatedAt:       now.Add(-72 * time.Hour),
		},
	}

	return exercises, nil
}
