package api

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/ogen-go/ogen/middleware"

	"github.com/zorcal/sbgfit/backend/api/internal/openapi"
	"github.com/zorcal/sbgfit/backend/pkg/slogctx"
	"github.com/zorcal/sbgfit/backend/pkg/tracectx"
)

func panicRecoveryMiddleware(log *slog.Logger) openapi.Middleware {
	return func(req middleware.Request, next middleware.Next) (resp middleware.Response, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := debug.Stack()

				log.ErrorContext(req.Context, "Panic recovered",
					"panic", rec,
					"stack", string(stack),
					"operation", req.OperationName,
					"operation_id", req.OperationID,
					"method", req.Raw.Method,
					"path", req.Raw.URL.Path,
				)

				err = fmt.Errorf("PANIC: %v", rec)
			}
		}()

		return next(req)
	}
}

func httpTraceMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			traceparent := cmp.Or(r.Header.Get("Traceparent"), uuid.NewString())

			ctx = tracectx.Set(ctx, traceparent)
			ctx = slogctx.Attach(ctx, "traceparent", traceparent)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func httpLoggingMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			rr := &responseRecorder{ResponseWriter: w}

			defer func(ctx context.Context) {
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

				log.LogAttrs(ctx, slog.LevelInfo, "HTTP request completed", attrs...)
			}(r.Context())

			next.ServeHTTP(rr, r)
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter

	statusCode int
}

func (sr *responseRecorder) WriteHeader(code int) {
	if sr.statusCode == 0 {
		sr.statusCode = code
	}
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *responseRecorder) Write(data []byte) (int, error) {
	if sr.statusCode == 0 {
		sr.statusCode = http.StatusOK
	}
	return sr.ResponseWriter.Write(data)
}

// wrapHTTPMiddleware wraps the given handler with the provided middlewares.
// Middlewares are applied in the order they are provided, meaning the first
// middleware in the slice will be the outermost wrapper.
func wrapHTTPMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		if middleware != nil {
			handler = middleware(handler)
		}
	}
	return handler
}
