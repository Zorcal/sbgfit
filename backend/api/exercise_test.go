package api_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/zorcal/sbgfit/backend/api"
	"github.com/zorcal/sbgfit/backend/api/internal/openapi"
	"github.com/zorcal/sbgfit/backend/internal/core/mdl"
	"github.com/zorcal/sbgfit/backend/internal/testingx"
	"github.com/zorcal/sbgfit/backend/pkg/ptr"
)

func TestGetExercises(t *testing.T) {
	now := time.Now()

	exerciseID1 := uuid.New()
	exerciseID2 := uuid.New()

	exerciseSvc := &MockedExerciseServiced{
		ExercisesFunc: func(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error) {
			exs := []mdl.Exercise{
				{
					ID:             exerciseID1,
					Name:           "Push Up",
					Category:       "strength",
					Description:    ptr.To("A bodyweight exercise targeting chest and triceps"),
					Instructions:   []string{"Place hands on ground", "Lower body", "Push up"},
					EquipmentTypes: []string{"bodyweight"},
					PrimaryMuscles: []string{"chest", "triceps"},
					Tags:           []string{"beginner-friendly", "functional"},
					CreatedAt:      now.AddDate(0, -2, 0),
					UpdatedAt:      now.AddDate(0, -1, 0),
				},
				{
					ID:             exerciseID2,
					Name:           "Pull Up",
					Category:       "strength",
					Description:    ptr.To("An upper body exercise targeting back and biceps"),
					Instructions:   []string{"Hang from bar", "Pull body up", "Lower with control"},
					EquipmentTypes: []string{"barbell"},
					PrimaryMuscles: []string{"back", "biceps"},
					Tags:           []string{"advanced", "strength-endurance"},
					CreatedAt:      now.AddDate(0, -1, 0),
					UpdatedAt:      now.AddDate(0, 0, -7),
				},
			}
			return exs, nil
		},
	}

	cfg := api.Config{
		Log:             testingx.NewLogger(t),
		ExerciseService: exerciseSvc,
	}

	srv := testServer(t, cfg)

	resp := makeRequest(t, srv, http.MethodGet, "/api/v1/exercises", nil)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusOK)
	}

	gotResp := testingx.DecodeJSON[openapi.ExerciseResponse](t, resp.Body)

	wantResp := openapi.ExerciseResponse{
		Data: []openapi.Exercise{
			{
				ID:       exerciseID1,
				Name:     "Push Up",
				Category: "strength",
				Description: openapi.OptNilString{
					Value: "A bodyweight exercise targeting chest and triceps",
					Set:   true,
				},
				Instructions: []string{"Place hands on ground", "Lower body", "Push up"},
				EquipmentTypes: []openapi.EquipmentType{
					openapi.EquipmentType("bodyweight"),
				},
				PrimaryMuscles: []openapi.PrimaryMuscle{
					openapi.PrimaryMuscle("chest"),
					openapi.PrimaryMuscle("triceps"),
				},
				Tags: []openapi.ExerciseTag{
					openapi.ExerciseTag("beginner-friendly"),
					openapi.ExerciseTag("functional"),
				},
				CreatedAt: now.AddDate(0, -2, 0),
				UpdatedAt: now.AddDate(0, -1, 0),
			},
			{
				ID:       exerciseID2,
				Name:     "Pull Up",
				Category: "strength",
				Description: openapi.OptNilString{
					Value: "An upper body exercise targeting back and biceps",
					Set:   true,
				},
				Instructions: []string{"Hang from bar", "Pull body up", "Lower with control"},
				EquipmentTypes: []openapi.EquipmentType{
					openapi.EquipmentType("barbell"),
				},
				PrimaryMuscles: []openapi.PrimaryMuscle{
					openapi.PrimaryMuscle("back"),
					openapi.PrimaryMuscle("biceps"),
				},
				Tags: []openapi.ExerciseTag{
					openapi.ExerciseTag("advanced"),
					openapi.ExerciseTag("strength-endurance"),
				},
				CreatedAt: now.AddDate(0, -1, 0),
				UpdatedAt: now.AddDate(0, 0, -7),
			},
		},
	}

	testingx.AssertDiff(t, gotResp, wantResp, cmpopts.EquateApproxTime(time.Second))
}

func TestGetExercises_error(t *testing.T) {
	exerciseSvc := &MockedExerciseServiced{
		ExercisesFunc: func(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error) {
			return nil, errors.New("some error")
		},
	}

	cfg := api.Config{
		Log:             testingx.NewLogger(t),
		ExerciseService: exerciseSvc,
	}

	srv := testServer(t, cfg)

	resp := makeRequest(t, srv, http.MethodGet, "/api/v1/exercises", nil)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	gotResp := testingx.DecodeJSON[openapi.ErrorResponse](t, resp.Body)

	wantResp := openapi.ErrorResponse{
		Error: "Internal Server Error",
	}

	testingx.AssertDiff(t, gotResp, wantResp)
}

func TestGetExercises_queryParams(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		wantFilter  mdl.ExerciseFilter
	}{
		{
			name:        "no filters",
			queryParams: "",
			wantFilter:  mdl.ExerciseFilter{},
		},
		{
			name:        "name",
			queryParams: "?name=Push+Up",
			wantFilter: mdl.ExerciseFilter{
				Name: ptr.To("Push Up"),
			},
		},
		{
			name:        "category",
			queryParams: "?category=strength",
			wantFilter: mdl.ExerciseFilter{
				Category: ptr.To("strength"),
			},
		},
		{
			name:        "equipment types",
			queryParams: "?equipmentTypes=barbell",
			wantFilter: mdl.ExerciseFilter{
				EquipmentTypes: []string{"barbell"},
			},
		},
		{
			name:        "single tag",
			queryParams: "?tags=crossfit",
			wantFilter: mdl.ExerciseFilter{
				Tags: []string{"crossfit"},
			},
		},
		{
			name:        "multiple tags",
			queryParams: "?tags=crossfit,advanced",
			wantFilter: mdl.ExerciseFilter{
				Tags: []string{"crossfit", "advanced"},
			},
		},
		{
			name:        "multiple filters",
			queryParams: "?name=Deadlift&category=strength&equipmentTypes=barbell",
			wantFilter: mdl.ExerciseFilter{
				Name:           ptr.To("Deadlift"),
				Category:       ptr.To("strength"),
				EquipmentTypes: []string{"barbell"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exerciseSvc := &MockedExerciseServiced{
				ExercisesFunc: func(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error) {
					testingx.AssertDiff(t, fltr, tt.wantFilter)
					return nil, nil
				},
			}

			cfg := api.Config{
				Log:             testingx.NewLogger(t),
				ExerciseService: exerciseSvc,
			}

			srv := testServer(t, cfg)

			resp := makeRequest(t, srv, http.MethodGet, "/api/v1/exercises"+tt.queryParams, nil)

			if resp.StatusCode != http.StatusOK {
				t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusOK)
			}
		})
	}
}

func TestGetExercises_queryParams_invalid(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		wantError   string
	}{
		{
			name:        "category",
			queryParams: "?category=invalid_category",
			wantError:   "operation GetExercises: decode params: query: \"category\": invalid value: invalid_category",
		},
		{
			name:        "equipmentTypes",
			queryParams: "?equipmentTypes=barbell,invalid_equipment,dumbbells",
			wantError:   `operation GetExercises: decode params: query: "equipmentTypes": invalid: [1] (invalid value: invalid_equipment)`,
		},
		{
			name:        "primaryMuscles",
			queryParams: "?primaryMuscles=chest,invalid_muscle",
			wantError:   `operation GetExercises: decode params: query: "primaryMuscles": invalid: [1] (invalid value: invalid_muscle)`,
		},
		{
			name:        "tags",
			queryParams: "?tags=crossfit,invalid_tag",
			wantError:   `operation GetExercises: decode params: query: "tags": invalid: [1] (invalid value: invalid_tag)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exerciseSvc := &MockedExerciseServiced{
				ExercisesFunc: func(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error) {
					return []mdl.Exercise{}, nil
				},
			}

			cfg := api.Config{
				Log:             testingx.NewLogger(t),
				ExerciseService: exerciseSvc,
			}

			srv := testServer(t, cfg)

			resp := makeRequest(t, srv, http.MethodGet, "/api/v1/exercises"+tt.queryParams, nil)

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusBadRequest)
			}

			gotResp := testingx.DecodeJSON[openapi.ErrorResponse](t, resp.Body)

			wantResp := openapi.ErrorResponse{
				Error: tt.wantError,
			}

			testingx.AssertDiff(t, gotResp, wantResp)
		})
	}
}
