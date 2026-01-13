package exercise

import (
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/zorcal/sbgfit/backend/internal/core/mdl"
)

func exercisesQuery(fltr mdl.ExerciseFilter, limit, offset int) (string, pgx.NamedArgs) {
	var q strings.Builder

	q.WriteString(`
		SELECT
			*,
			COUNT(*) OVER() as total_count
		FROM (
			SELECT
				e.external_id,
				e.name,
				c.code as category_code,
				e.description,
				e.instructions,
				COALESCE(
					ARRAY_AGG(DISTINCT et.code) FILTER (WHERE et.code IS NOT NULL),
					ARRAY[]::text[]
				) as equipment_types,
				COALESCE(
					ARRAY_AGG(DISTINCT pm.code) FILTER (WHERE pm.code IS NOT NULL),
					ARRAY[]::text[]
				) as primary_muscles,
				COALESCE(
					ARRAY_AGG(DISTINCT tag.code) FILTER (WHERE tag.code IS NOT NULL),
					ARRAY[]::text[]
				) as tags,
				e.created_at,
				e.updated_at
			FROM sbgfit.exercises e
			LEFT JOIN sbgfit.exercise_categories c ON e.category_id = c.id
			LEFT JOIN sbgfit.exercise_equipment ee ON e.id = ee.exercise_id
			LEFT JOIN sbgfit.equipment_types et ON ee.equipment_type_id = et.id
			LEFT JOIN sbgfit.exercise_primary_muscles epm ON e.id = epm.exercise_id
			LEFT JOIN sbgfit.primary_muscles pm ON epm.primary_muscle_id = pm.id
			LEFT JOIN sbgfit.exercise_exercise_tags eet ON e.id = eet.exercise_id
			LEFT JOIN sbgfit.exercise_tags tag ON eet.exercise_tag_id = tag.id
			GROUP BY e.id, e.external_id, e.name, c.code, e.description, e.instructions, e.created_at, e.updated_at
		) AS exercise_data`)

	args := make(pgx.NamedArgs)

	var predicates []string

	if fltr.Name != nil {
		predicates = append(predicates, "exercise_data.name ILIKE @name")
		args["name"] = "%" + *fltr.Name + "%"
	}

	if fltr.Category != nil {
		predicates = append(predicates, "exercise_data.category_code = @category")
		args["category"] = *fltr.Category
	}

	if len(fltr.EquipmentTypes) > 0 {
		predicates = append(predicates, "exercise_data.equipment_types && @equipmentTypes")
		args["equipmentTypes"] = fltr.EquipmentTypes
	}

	if len(fltr.PrimaryMuscles) > 0 {
		predicates = append(predicates, "exercise_data.primary_muscles && @primaryMuscles")
		args["primaryMuscles"] = fltr.PrimaryMuscles
	}

	if len(fltr.Tags) > 0 {
		predicates = append(predicates, "exercise_data.tags && @tags")
		args["tags"] = fltr.Tags
	}

	if len(predicates) > 0 {
		q.WriteString(" WHERE ")
		q.WriteString(strings.Join(predicates, " AND "))
	}

	q.WriteString(`
		ORDER BY name COLLATE natsort
		LIMIT @limit OFFSET @offset`)

	args["limit"] = limit
	args["offset"] = offset

	return q.String(), args
}
