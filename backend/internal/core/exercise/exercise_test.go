package exercise

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/zorcal/sbgfit/backend/internal/core/mdl"
	"github.com/zorcal/sbgfit/backend/internal/data/pgtest"
	"github.com/zorcal/sbgfit/backend/internal/testingx"
	"github.com/zorcal/sbgfit/backend/pkg/ptr"
)

func TestExercises(t *testing.T) {
	ctx := context.Background()

	pool := pgtest.NewWithSeed(t, ctx)

	svc := NewService(pool)

	airSquats := mdl.Exercise{
		Name:           "Air Squats",
		Category:       "strength",
		Description:    ptr.To("Bodyweight squat focusing on proper hip and knee movement"),
		Instructions:   []string{"Stand with feet shoulder-width", "Lower hips back and down", "Keep chest up", "Drive through heels", "Return to standing"},
		EquipmentTypes: []string{"bodyweight"},
		PrimaryMuscles: []string{"glutes", "legs"},
		Tags:           []string{"beginner-friendly", "crossfit", "functional"},
	}

	assaultBike := mdl.Exercise{
		Name:           "Assault Bike",
		Category:       "cardio",
		Description:    ptr.To("High-intensity cardio using air resistance bike with moving handles"),
		Instructions:   []string{"Sit on bike with feet on pedals", "Grip moving handles", "Push and pull with arms", "Pedal with legs simultaneously", "Maintain steady breathing"},
		EquipmentTypes: []string{"assault-bike"},
		PrimaryMuscles: []string{"core", "full-body", "legs"},
		Tags:           []string{"advanced", "conditioning", "crossfit", "hyrox"},
	}

	barbellBackSquat := mdl.Exercise{
		Name:           "Barbell Back Squat",
		Category:       "strength",
		Description:    ptr.To("Fundamental squatting movement with barbell on back"),
		Instructions:   []string{"Position barbell on upper back", "Stand with feet shoulder-width", "Descend by sitting back", "Drive through heels to stand", "Keep chest up throughout"},
		EquipmentTypes: []string{"barbell"},
		PrimaryMuscles: []string{"core", "glutes", "legs"},
		Tags:           []string{"crossfit", "functional", "strength-endurance"},
	}

	barbellBenchPress := mdl.Exercise{
		Name:           "Barbell Bench Press",
		Category:       "strength",
		Description:    ptr.To("Classic upper body pressing movement with barbell"),
		Instructions:   []string{"Lie on bench with barbell racked", "Grip barbell slightly wider than shoulders", "Lower bar to chest with control", "Press bar straight up", "Lock out arms at top"},
		EquipmentTypes: []string{"barbell"},
		PrimaryMuscles: []string{"chest", "shoulders", "triceps"},
		Tags:           []string{"functional", "strength-endurance"},
	}

	barbellBentOverRows := mdl.Exercise{
		Name:           "Barbell Bent-over Rows",
		Category:       "strength",
		Description:    ptr.To("Pulling movement with barbell targeting back muscles"),
		Instructions:   []string{"Hinge at hips holding barbell", "Keep back straight and core tight", "Pull barbell to lower chest", "Squeeze shoulder blades together", "Lower with control"},
		EquipmentTypes: []string{"barbell"},
		PrimaryMuscles: []string{"back", "biceps", "core"},
		Tags:           []string{"functional", "strength-endurance"},
	}

	burpees := mdl.Exercise{
		Name:           "Burpees",
		Category:       "cardio",
		Description:    ptr.To("From standing, squat down, jump back to plank, do a push-up, jump feet back to squat, then jump up with arms overhead"),
		Instructions:   []string{"Start standing", "Squat down hands on ground", "Jump back to plank", "Do push-up", "Jump feet to squat", "Jump up arms overhead"},
		EquipmentTypes: []string{"bodyweight"},
		PrimaryMuscles: []string{"full-body"},
		Tags:           []string{"competition", "conditioning", "crossfit", "functional", "hyrox"},
	}

	dips := mdl.Exercise{
		Name:           "Dips",
		Category:       "strength",
		Description:    ptr.To("Bodyweight exercise targeting triceps and chest"),
		Instructions:   []string{"Support body on parallel bars", "Lower body down", "Push back to start", "Keep body upright", "Control the movement"},
		EquipmentTypes: []string{"bodyweight"},
		PrimaryMuscles: []string{"chest", "shoulders", "triceps"},
		Tags:           []string{"functional", "strength-endurance"},
	}

	tests := []struct {
		name           string
		fltr           mdl.ExerciseFilter
		pageSize       int
		pageNumber     int
		want           []mdl.Exercise
		wantTotalCount int
	}{
		{
			name:           "no filters",
			fltr:           mdl.ExerciseFilter{},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{airSquats, assaultBike},
			wantTotalCount: 35,
		},
		{
			name:           "filter by name",
			fltr:           mdl.ExerciseFilter{Name: ptr.To("aIr")},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{airSquats},
			wantTotalCount: 1,
		},
		{
			name:           "filter by category",
			fltr:           mdl.ExerciseFilter{Category: ptr.To("cardio")},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{assaultBike, burpees},
			wantTotalCount: 7,
		},
		{
			name:           "filter by equipment types",
			fltr:           mdl.ExerciseFilter{EquipmentTypes: []string{"barbell"}},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{barbellBackSquat, barbellBenchPress},
			wantTotalCount: 8,
		},
		{
			name:           "filter by multiple equipment types",
			fltr:           mdl.ExerciseFilter{EquipmentTypes: []string{"bodyweight", "kettlebell"}},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{airSquats, burpees},
			wantTotalCount: 11,
		},
		{
			name:           "filter by primary muscles",
			fltr:           mdl.ExerciseFilter{PrimaryMuscles: []string{"chest"}},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{barbellBenchPress, dips},
			wantTotalCount: 4,
		},
		{
			name:           "filter by multiple primary muscles",
			fltr:           mdl.ExerciseFilter{PrimaryMuscles: []string{"core", "full-body"}},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{assaultBike, barbellBackSquat},
			wantTotalCount: 24,
		},
		{
			name:           "filter by tags",
			fltr:           mdl.ExerciseFilter{Tags: []string{"crossfit"}},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{airSquats, assaultBike},
			wantTotalCount: 22,
		},
		{
			name:           "filter by multiple tags",
			fltr:           mdl.ExerciseFilter{Tags: []string{"functional", "strength-endurance"}},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{airSquats, barbellBackSquat},
			wantTotalCount: 27,
		},
		{
			name: "multiple filters",
			fltr: mdl.ExerciseFilter{
				Category:       ptr.To("strength"),
				EquipmentTypes: []string{"barbell"},
				Tags:           []string{"functional"},
			},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{barbellBackSquat, barbellBenchPress},
			wantTotalCount: 7,
		},
		{
			name:           "pagination - first page",
			fltr:           mdl.ExerciseFilter{},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{airSquats, assaultBike},
			wantTotalCount: 35,
		},
		{
			name:           "pagination - second page",
			fltr:           mdl.ExerciseFilter{},
			pageSize:       2,
			pageNumber:     2,
			want:           []mdl.Exercise{barbellBackSquat, barbellBenchPress},
			wantTotalCount: 35,
		},
		{
			name:           "pagination with filters - first page",
			fltr:           mdl.ExerciseFilter{Category: ptr.To("strength")},
			pageSize:       2,
			pageNumber:     1,
			want:           []mdl.Exercise{airSquats, barbellBackSquat},
			wantTotalCount: 27,
		},
		{
			name:           "pagination with filters - second page",
			fltr:           mdl.ExerciseFilter{Category: ptr.To("strength")},
			pageSize:       2,
			pageNumber:     2,
			want:           []mdl.Exercise{barbellBenchPress, barbellBentOverRows},
			wantTotalCount: 27,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotTotalCount, err := svc.Exercises(ctx, tt.fltr, tt.pageSize, tt.pageNumber)
			if err != nil {
				t.Fatalf("Exercises(%+v) error = %v, want no error", tt.fltr, err)
			}

			if gotTotalCount != tt.wantTotalCount {
				t.Errorf("Exercises(%+v) total count = %d, want %d", tt.fltr, gotTotalCount, tt.wantTotalCount)
			}

			diffOpts := cmp.Options{
				cmpopts.IgnoreFields(mdl.Exercise{}, "ID", "CreatedAt", "UpdatedAt"), // Ignore generated fields
			}
			testingx.AssertDiff(t, got, tt.want, diffOpts)
		})
	}
}
