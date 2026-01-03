package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ardanlabs/conf/v3"
)

type Config struct {
	conf.Version

	Environment string `conf:"default:local"`
	Web         struct {
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:10s"`
		IdleTimeout     time.Duration `conf:"default:120s"`
		ShutdownTimeout time.Duration `conf:"default:20s"`
		Addr            string        `conf:"default:127.0.0.1:4250"`
	}
}

// LogValue implements slog.LogValuer.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("environment", c.Environment),
		slog.String("build", c.Build),
		slog.String("desc", c.Desc),
		slog.Group("web",
			slog.String("read_timeout", c.Web.ReadTimeout.String()),
			slog.String("write_timeout", c.Web.WriteTimeout.String()),
			slog.String("idle_timeout", c.Web.IdleTimeout.String()),
			slog.String("shutdown_timeout", c.Web.ShutdownTimeout.String()),
			slog.String("addr", c.Web.Addr),
		),
	)
}

func (c Config) Validate() error {
	if c.Environment != "local" && c.Environment != "prod" {
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	return nil
}
