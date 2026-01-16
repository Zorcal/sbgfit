package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"

	"github.com/zorcal/sbgfit/backend/api"
	"github.com/zorcal/sbgfit/backend/internal/core/exercise"
	"github.com/zorcal/sbgfit/backend/internal/data/pgdb"
	"github.com/zorcal/sbgfit/backend/internal/data/schema"
	"github.com/zorcal/sbgfit/backend/internal/telemetry"
	"github.com/zorcal/sbgfit/backend/pkg/slogctx"
)

// appVersion should be set at build time using -ldflags.
var appVersion = "dev"

func main() {
	ctx := context.Background()

	cfg := Config{
		Version: conf.Version{
			Build: appVersion,
			Desc:  "Fitness application.",
		},
	}
	if help, err := conf.Parse("SBGFIT", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Fprint(os.Stdout, help)
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
		os.Exit(1)
	}

	log := slog.New(logHandler(cfg.Environment))

	if err := run(ctx, cfg, log); err != nil {
		log.ErrorContext(ctx, "Run error", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg Config, log *slog.Logger) (retErr error) {
	log.InfoContext(ctx, "Starting...", "config", cfg)

	telemetryConfig := telemetry.Config{
		Enabled:  cfg.Telemetry.Enabled,
		Endpoint: cfg.Telemetry.Endpoint,
		Insecure: cfg.Telemetry.Insecure,
	}
	cleanupTracing, err := telemetry.InitTracing(ctx, "sbgfit-backend", appVersion, telemetryConfig, log)
	if err != nil {
		return fmt.Errorf("initialize tracing: %w", err)
	}
	defer cleanupTracing()

	connStr := pgdb.ConnStr(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLEnabled)

	if err := schema.Migrate(ctx, connStr); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("parse database pool config: %w", err)
	}

	pool, err := pgdb.NewPool(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("new database pool: %w", err)
	}
	defer pool.Close()

	if err := pgdb.StatusCheck(ctx, pool); err != nil {
		return fmt.Errorf("status check database connection: %w", err)
	}

	if err := schema.SeedData(ctx, pool); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	exerciseSvc := exercise.NewService(pool)

	handler, err := api.NewHandler(api.Config{
		Log:             log,
		ExerciseService: exerciseSvc,
	})
	if err != nil {
		return fmt.Errorf("create handler: %w", err)
	}

	srv := http.Server{
		Addr:         cfg.Web.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     slog.NewLogLogger(log.Handler(), slog.LevelInfo),
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	srvErrs := make(chan error, 1)

	go func() {
		log.InfoContext(ctx, "HTTP server started", "host", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			srvErrs <- fmt.Errorf("listen and serve: %w", err)
		}
	}()

	defer log.InfoContext(ctx, "HTTP server stopped")

	select {
	case err := <-srvErrs:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.InfoContext(ctx, "Graceful shutdown started", "signal", sig)
		defer log.InfoContext(ctx, "Shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return fmt.Errorf("could not stop HTTP server gracefully: %w", err)
		}
	}

	return nil
}

func logHandler(env string) slog.Handler {
	var h slog.Handler
	if env == "local" {
		h = tint.NewHandler(os.Stdout, nil)
	} else {
		h = slog.NewJSONHandler(os.Stdout, nil)
	}
	h = slogctx.NewHandler(h)
	return h
}
