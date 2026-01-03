package handler

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"

	"github.com/zorcal/sbgfit/backend/pkg/httprouter"
	"github.com/zorcal/sbgfit/backend/pkg/slogctx"
	"github.com/zorcal/sbgfit/backend/pkg/tracectx"
)

func traceMiddleware() httprouter.Middleware {
	return func(next httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			traceparent := cmp.Or(r.Header.Get("Traceparent"), uuid.NewString())

			ctx = tracectx.Set(ctx, traceparent)
			ctx = slogctx.Attach(ctx, "traceparent", traceparent)

			return next(w, r.WithContext(ctx))
		}
	}
}

func loggingMiddleware(log *slog.Logger) httprouter.Middleware {
	return func(next httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) (retErr error) {
			now := time.Now()

			rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			defer func() {
				attrs := []slog.Attr{
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("x_forwarded_for", r.Header.Get("X-Forwarded-For")),
					slog.Int("status_code", rr.statusCode),
					slog.String("user_agent", r.Header.Get("User-Agent")),
					slog.Time("started_at", now),
					slog.Time("finished_at", time.Now()),
					slog.Int64("took_ms", time.Since(now).Milliseconds()),
				}
				if retErr != nil {
					// Should not happen. Caught by error middleware.
					attrs = append(attrs, slog.String("error", retErr.Error()))
				}
				log.LogAttrs(r.Context(), logLevel(rr.statusCode), "HTTP request completed", attrs...)
			}()

			return next(rr, r)
		}
	}
}

// responseRecorder is a wrapper around http.ResponseWriter to capture the
// status code.
type responseRecorder struct {
	http.ResponseWriter

	statusCode int
}

func (rw *responseRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func logLevel(statusCode int) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError
	case statusCode >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func panicRecovery(log *slog.Logger) httprouter.Middleware {
	return func(next httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()

					log.ErrorContext(r.Context(), "Panic recovered",
						"panic", rec,
						"stack", string(stack),
						"method", r.Method,
						"path", r.URL.Path,
					)

					err = fmt.Errorf("PANIC: %v", rec)
				}
			}()

			return next(w, r)
		}
	}
}

func errorMiddleware(log *slog.Logger) httprouter.Middleware {
	return func(next httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			err := next(w, r)
			if err == nil {
				return nil
			}

			ctx := r.Context()

			var httpErr *httpError
			if !errors.As(err, &httpErr) {
				httpErr = &httpError{
					StatusCode:      http.StatusInternalServerError,
					ExternalMessage: http.StatusText(http.StatusInternalServerError),
					InternalErr:     err,
				}
			}

			log.ErrorContext(ctx, "Request error", "error", httpErr)

			w.WriteHeader(httpErr.StatusCode)

			if err := json.NewEncoder(w).Encode(httpErr); err != nil {
				log.ErrorContext(ctx, "Failed to encode error response", "error", err)
			}

			return nil
		}
	}
}
