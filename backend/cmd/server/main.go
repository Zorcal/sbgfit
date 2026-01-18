package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ardanlabs/conf/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"github.com/pgx-contrib/pgxotel"

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

	dbConnStr := pgdb.ConnStr(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLEnabled)

	if err := schema.Migrate(ctx, dbConnStr); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	poolQueryParams := url.Values{}
	if cfg.DB.Pool.MaxConns > 0 {
		poolQueryParams.Set("pool_max_conns", strconv.Itoa(cfg.DB.Pool.MaxConns))
	}
	if cfg.DB.Pool.MinConns > 0 {
		poolQueryParams.Set("pool_min_conns", strconv.Itoa(cfg.DB.Pool.MinConns))
	}
	if cfg.DB.Pool.MaxConnLifetime > 0 {
		poolQueryParams.Set("pool_max_conn_lifetime", cfg.DB.Pool.MaxConnLifetime.String())
	}
	if cfg.DB.Pool.MaxConnIdleTime > 0 {
		poolQueryParams.Set("pool_max_conn_idle_time", cfg.DB.Pool.MaxConnIdleTime.String())
	}
	if cfg.DB.Pool.HealthCheckPeriod > 0 {
		poolQueryParams.Set("pool_health_check_period", cfg.DB.Pool.HealthCheckPeriod.String())
	}
	if cfg.DB.Pool.MaxConnLifetimeJitter > 0 {
		poolQueryParams.Set("pool_max_conn_lifetime_jitter", cfg.DB.Pool.MaxConnLifetimeJitter.String())
	}

	poolConnstr, err := mergeURLQueryParams(dbConnStr, poolQueryParams)
	if err != nil {
		return fmt.Errorf("merge pool query params with db connstr: %w", err)
	}

	poolCfg, err := pgxpool.ParseConfig(poolConnstr)
	if err != nil {
		return fmt.Errorf("parse database pool config: %w", err)
	}
	poolCfg.ConnConfig.Tracer = &pgxotel.QueryTracer{
		Name: "sbgfit-postgres",
	}

	pool, err := pgdb.NewPool(ctx, poolCfg)
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

func mergeURLQueryParams(rawURL string, newParams url.Values) (string, error) {
	if len(newParams) == 0 {
		return rawURL, nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse raw URL: %w", err)
	}

	q := parsed.Query()
	for key, values := range newParams {
		for _, v := range values {
			q.Set(key, v)
		}
	}

	parsed.RawQuery = q.Encode()

	return parsed.String(), nil
}
