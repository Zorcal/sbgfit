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
	"github.com/lmittmann/tint"
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

	srv := http.Server{
		Addr: cfg.Web.Addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello world!")
		}),
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
