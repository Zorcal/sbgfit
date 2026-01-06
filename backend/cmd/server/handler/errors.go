package handler

import (
	"fmt"
	"log/slog"
)

type httpError struct {
	StatusCode      int    `json:"-"`
	ExternalMessage string `json:"error,omitzero"`
	InternalErr     error  `json:"-"`
}

func (e *httpError) Error() string {
	if e.InternalErr != nil {
		return fmt.Sprintf("%s: %v", e.ExternalMessage, e.InternalErr)
	}
	return e.ExternalMessage
}

func (e *httpError) Unwrap() error {
	return e.InternalErr
}

// LogValue implements slog.LogValuer.
func (e *httpError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("status_code", e.StatusCode),
		slog.String("external", e.ExternalMessage),
		slog.String("inernal", e.InternalErr.Error()),
	)
}
