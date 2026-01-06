package handler

import (
	"github.com/zorcal/sbgfit/backend/pkg/httprouter"
)

func routes(r *httprouter.Router, cfg Config) {
	r.Handle("GET /api/v1/exercises", getExercisesHandler(cfg.ExerciseService))
}
