package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/zorcal/sbgfit/backend/api/internal/openapi"
)

type Config struct {
	Log             *slog.Logger
	ExerciseService ExerciseService
}

func NewHandler(cfg Config) (http.Handler, error) {
	mux := http.NewServeMux()

	v1Handler, err := newV1Handler(cfg)
	if err != nil {
		return nil, fmt.Errorf("create v1 handler: %w", err)
	}

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Handler))

	return mux, nil
}

func newV1Handler(cfg Config) (http.Handler, error) {
	srv, err := openapi.NewServer(
		&api{
			log:         cfg.Log,
			exerciseSvc: cfg.ExerciseService,
		},
		openapi.WithMiddleware(middleware.ChainMiddlewares(
			panicRecoveryMiddleware(cfg.Log),
		)),
		openapi.WithErrorHandler(errorHandler),
	)
	if err != nil {
		return nil, fmt.Errorf("create openapi server: %w", err)
	}

	// TODO: make global to mux
	h := wrapHTTPMiddleware(srv,
		httpTraceMiddleware(),
		httpLoggingMiddleware(cfg.Log),
	)

	return h, nil
}

func errorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("error")
	e.StrEscape(err.Error())
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
	jx.PutEncoder(e)
}
