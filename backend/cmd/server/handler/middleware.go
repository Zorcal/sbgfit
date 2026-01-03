package handler

import (
	"cmp"
	"log/slog"
	"net/http"
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
