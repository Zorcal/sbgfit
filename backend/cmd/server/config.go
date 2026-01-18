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
	DB struct {
		User       string `conf:"default:postgres"`
		Password   string `conf:"default:postgres,mask"`
		Host       string `conf:"default:sbgfit-postgres"`
		Port       int    `conf:"default:5433"`
		Name       string `conf:"default:sbgfit"`
		SSLEnabled bool   `conf:"default:false"`
		Pool       struct {
			MaxConns              int
			MinConns              int
			MaxConnLifetime       time.Duration `conf:"default:1h"`
			MaxConnIdleTime       time.Duration `conf:"default:30m"`
			HealthCheckPeriod     time.Duration `conf:"default:60s"`
			MaxConnLifetimeJitter time.Duration `conf:"default:0s"`
		}
	}
	Telemetry struct {
		Enabled  bool   `conf:"default:true"`
		Endpoint string `conf:"default:127.0.0.1:4317"`
		Insecure bool   `conf:"default:true"`
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
		slog.Group("db",
			slog.String("user", c.DB.User),
			slog.String("host", c.DB.Host),
			slog.Int("port", c.DB.Port),
			slog.String("name", c.DB.Name),
			slog.Bool("ssl_enabled", c.DB.SSLEnabled),
			slog.Group("pool",
				slog.Int("max_conns", c.DB.Pool.MaxConns),
				slog.Int("min_conns", c.DB.Pool.MinConns),
				slog.Duration("max_conn_lifetime", c.DB.Pool.MaxConnLifetime),
				slog.Duration("max_conn_idle_time", c.DB.Pool.MaxConnIdleTime),
				slog.Duration("health_check_period", c.DB.Pool.HealthCheckPeriod),
				slog.Duration("max_conn_lifetime_jitter", c.DB.Pool.MaxConnLifetimeJitter),
			),
		),
		slog.Group("telemetry",
			slog.Bool("enabled", c.Telemetry.Enabled),
			slog.String("endpoint", c.Telemetry.Endpoint),
			slog.Bool("insecure", c.Telemetry.Insecure),
		),
	)
}

func (c Config) Validate() error {
	if c.Environment != "local" && c.Environment != "prod" {
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	return nil
}
