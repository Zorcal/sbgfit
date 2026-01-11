package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/zorcal/sbgfit/backend/core/mdl"
	"github.com/zorcal/sbgfit/backend/core/testingx"
	"github.com/zorcal/sbgfit/backend/pkg/ptr"
)

func TestGetExercisesHandler(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		mockResponse []mdl.Exercise
	}{
		{
			name:        "without filters",
			queryParams: "",
			mockResponse: []mdl.Exercise{
				{
					ID:              uuid.New(),
					Name:            "Push Up",
					Category:        "strength",
					EquipmentTypes:  []string{"bodyweight"},
					PrimaryMuscles:  []string{"chest"},
					Tags:            []string{"functional"},
					CreatedByUserID: uuid.New(),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
				{
					ID:              uuid.New(),
					Name:            "Pull Up",
					Category:        "strength",
					EquipmentTypes:  []string{"bodyweight"},
					PrimaryMuscles:  []string{"back"},
					Tags:            []string{"functional"},
					CreatedByUserID: uuid.New(),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
		},
		{
			name:        "with filters",
			queryParams: "?name=squat&category=strength&equipmentTypes=barbell,dumbbells&primaryMuscles=quads,glutes&tags=crossfit,functional&createdByUser=true",
			mockResponse: []mdl.Exercise{
				{
					ID:              uuid.New(),
					Name:            "Back Squat",
					Category:        "strength",
					EquipmentTypes:  []string{"barbell"},
					PrimaryMuscles:  []string{"quads", "glutes"},
					Tags:            []string{"crossfit", "functional"},
					CreatedByUserID: uuid.New(),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
		},
		{
			name:         "empty response",
			queryParams:  "?name=nonexistent",
			mockResponse: []mdl.Exercise{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockeExerciseServiced{
				ExercisesFunc: func(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error) {
					return tt.mockResponse, nil
				},
			}

			h := New(Config{
				Log:             testingx.NewLogger(t),
				ExerciseService: mockService,
			})
			srv := httptest.NewServer(h)

			client := srv.Client()

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, fmt.Sprintf("%s/api/v1/exercises%s", srv.URL, tt.queryParams), nil)
			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("error making request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("got status code = %d, want %d", resp.StatusCode, http.StatusOK)
			}

			var got response[[]mdl.Exercise]
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatalf("error decoding response body: %v", err)
			}

			want := response[[]mdl.Exercise]{
				Data: tt.mockResponse,
			}

			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("response diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestGetExercisesHandler_error(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockError      error
		wantStatusCode int
		wantResponse   map[string]any
	}{
		{
			name:           "service error",
			mockError:      errors.New("database connection failed"),
			wantStatusCode: http.StatusInternalServerError,
			wantResponse: map[string]any{
				"error": "Internal Server Error",
			},
		},
		{
			name:           "invalid query params",
			queryParams:    "?equipmentTypes=invalid_type",
			wantStatusCode: http.StatusBadRequest,
			wantResponse: map[string]any{
				"error": "Invalid filter query params",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockeExerciseServiced{
				ExercisesFunc: func(ctx context.Context, fltr mdl.ExerciseFilter) ([]mdl.Exercise, error) {
					return nil, tt.mockError
				},
			}

			h := New(Config{
				Log:             testingx.NewLogger(t),
				ExerciseService: mockService,
			})
			srv := httptest.NewServer(h)

			client := srv.Client()

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, fmt.Sprintf("%s/api/v1/exercises%s", srv.URL, tt.queryParams), nil)
			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("error making request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("got status code = %d, want %d", resp.StatusCode, tt.wantStatusCode)
			}

			var got map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatalf("error decoding response body: %v", err)
			}

			if diff := cmp.Diff(got, tt.wantResponse); diff != "" {
				t.Errorf("response diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestBuildExerciseFilter(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		want        mdl.ExerciseFilter
	}{
		{
			name:        "empty query parameters",
			queryParams: "",
			want:        mdl.ExerciseFilter{},
		},
		{
			name:        "all filters combined",
			queryParams: "?name=deadlift&category=strength&equipmentTypes=barbell&primaryMuscles=back,hamstrings&tags=functional,strength-endurance&createdByUser=true",
			want: mdl.ExerciseFilter{
				Name:           ptr.To("deadlift"),
				Category:       ptr.To("strength"),
				EquipmentTypes: []string{"barbell"},
				PrimaryMuscles: []string{"back", "hamstrings"},
				Tags:           []string{"functional", "strength-endurance"},
				CreatedByUser:  ptr.To(true),
			},
		},
		{
			name:        "createdByUser false",
			queryParams: "?createdByUser=false",
			want: mdl.ExerciseFilter{
				CreatedByUser: ptr.To(false),
			},
		},
		{
			name:        "empty string values ignored",
			queryParams: "?name=&category=&equipmentTypes=&primaryMuscles=&tags=&createdByUser=",
			want:        mdl.ExerciseFilter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/exercises"+tt.queryParams, nil)

			got, err := buildExerciseFilter(req)
			if err != nil {
				t.Fatalf("buildExerciseFilter(%q) error = %v, want nil", tt.queryParams, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("buildExerciseFilter(%q) diff (-want +got):\n%s", tt.queryParams, diff)
			}
		})
	}
}

func TestBuildExerciseFilter_error(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
	}{
		{
			name:        "invalid category",
			queryParams: "?category=invalid_category",
		},
		{
			name:        "invalid equipment type",
			queryParams: "?equipmentTypes=invalid_type",
		},
		{
			name:        "invalid primary muscle",
			queryParams: "?primaryMuscles=invalid_muscle",
		},
		{
			name:        "invalid exercise tag",
			queryParams: "?tags=invalid_tag",
		},
		{
			name:        "invalid createdByUser value",
			queryParams: "?createdByUser=maybe",
		},
		{
			name:        "mixed valid and invalid equipment types",
			queryParams: "?equipmentTypes=barbell,invalid_type,dumbbells",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/exercises"+tt.queryParams, nil)
			if _, err := buildExerciseFilter(req); err == nil {
				t.Fatalf("buildExerciseFilter(%q) error = nil, want error", tt.queryParams)
			}
		})
	}
}
