// Package exercise provides the application service for managing both the
// exercise library and user-created exercises. This includes retrieving
// predefined exercises from the library that users can browse and clone, as
// well as managing user-specific exercises that have been customized.
package exercise

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zorcal/sbgfit/backend/internal/core/mdl"
	"github.com/zorcal/sbgfit/backend/internal/telemetry"
)

// Service manages both the exercise library and user-created exercises. It
// provides read-only access to predefined exercises from the library that
// serve as templates, and will handle CRUD operations for user-specific
// exercises that have been cloned and customized.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a new exercise service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		pool: pool,
	}
}

// Exercises retrieves predefined exercises from the exercise library based on
// the provided filter criteria.
func (s *Service) Exercises(ctx context.Context, fltr mdl.ExerciseFilter, pageSize, pageNumber int) (exs []mdl.Exercise, totalCount int, retErr error) {
	ctx, span := telemetry.StartSpan(ctx, "exercise.Service.Exercises")
	defer span.End()

	offset := (pageNumber - 1) * pageSize

	exercisesQ, args := exercisesQuery(fltr, pageSize, offset)

	rows, err := s.pool.Query(ctx, exercisesQ, args)
	if err != nil {
		return nil, 0, fmt.Errorf("query exercises: %w", err)
	}
	defer rows.Close()

	type Result struct {
		dbExercise

		TotalCount int `db:"total_count"`
	}
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[Result])
	if err != nil {
		return nil, 0, fmt.Errorf("collect exercise rows: %w", err)
	}

	if len(result) > 0 {
		totalCount = result[0].TotalCount
	}

	exs = make([]mdl.Exercise, len(result))
	for i, row := range result {
		exs[i] = dbExerciseToModel(row.dbExercise)
	}

	return exs, totalCount, nil
}
