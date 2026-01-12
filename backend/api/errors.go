package api

import (
	"fmt"
	"log/slog"
)

type httpError struct {
	StatusCode      int
	ExternalMessage string
	InternalErr     error
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

func (e *httpError) LogValue() slog.Value {
	var internalErrStr string
	if e.InternalErr != nil {
		internalErrStr = e.InternalErr.Error()
	}
	return slog.GroupValue(
		slog.Int("status_code", e.StatusCode),
		slog.String("external", e.ExternalMessage),
		slog.String("internal", internalErrStr),
	)
}
