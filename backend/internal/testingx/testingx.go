// Package testingx extends the testing package from the standard library with
// application specific testing utilities.
package testingx

import (
	"log/slog"
	"testing"

	"github.com/lmittmann/tint"

	"github.com/zorcal/sbgfit/backend/pkg/slogctx"
)

// NewLogger creates a new logger for tests.
func NewLogger(t *testing.T) *slog.Logger {
	t.Helper()
	h := slogctx.NewHandler(tint.NewHandler(&Writer{t: t}, &tint.Options{Level: slog.LevelDebug}))
	return slog.New(h)
}

type Writer struct {
	t *testing.T
}

func (tw *Writer) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}
