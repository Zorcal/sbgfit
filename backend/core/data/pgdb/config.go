package pgdb

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config defines database and pool settings used to construct a pgx connection
// configuration. It exists to provide a safer and less error-prone alternative
// to manually assembling a PostgreSQL connection string.
//
// Zero values are treated as "unset" and omitted from the resulting connection.
// Pool-related fields map directly to pgxpool query parameters. The resulting
// connection URL enforces UTC timezone and sets sslmode based on DisableTLS.
type Config struct {
	// User is the database user.
	User string
	// Password is the database user password.
	Password string
	// Host is the database host. IPv6 addresses must be bracketed if a port is included.
	Host string
	// Name is the database name.
	Name string
	// Schema sets the PostgreSQL search_path if non-empty.
	Schema string
	// DisableTLS disables TLS by setting sslmode=disable. If false, sslmode=require is used.
	DisableTLS bool

	// PoolMaxConns is the maximum size of the pool. The default is the greater
	// of 4 or runtime.NumCPU().
	PoolMaxConns int
	// PoolMinConns is the minimum size of the pool. After a connection closes,
	// the pool may dip below PoolMinConns until the next health check.
	PoolMinConns int
	// PoolMaxConnLifetime is the maximum lifetime of a connection.
	PoolMaxConnLifetime time.Duration
	// PoolMaxConnIdleTime is the maximum idle time before a connection is closed.
	PoolMaxConnIdleTime time.Duration
	// PoolHealthCheckPeriod is the interval between pool health checks.
	PoolHealthCheckPeriod time.Duration
	// PoolMaxConnLifetimeJitter adds randomized delay after MaxConnLifetime
	// before closing a connection to avoid synchronized churn.
	PoolMaxConnLifetimeJitter time.Duration
}

// ParsePgxConfig validates the Config and parses it into a *pgxpool.Config
// using pgxpool.ParseConfig. The returned config is guaranteed to be valid
// for NewPool. Returns an error if validation fails or parsing
// fails.
func (cfg Config) ParsePgxConfig() (*pgxpool.Config, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	pgxCfg, err := pgxpool.ParseConfig(cfg.URL().String())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return pgxCfg, nil
}

// URL returns a PostgreSQL connection URL derived from the Config. The
// returned URL is suitable for pgxpool.ParseConfig and does not perform
// validation. The URL forces UTC timezone and sets sslmode based on
// DisableTLS.
func (cfg Config) URL() *url.URL {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")
	if cfg.Schema != "" {
		q.Set("search_path", cfg.Schema)
	}

	if cfg.PoolMaxConns > 0 {
		q.Set("pool_max_conns", strconv.Itoa(cfg.PoolMaxConns))
	}
	if cfg.PoolMinConns > 0 {
		q.Set("pool_min_conns", strconv.Itoa(cfg.PoolMinConns))
	}
	if cfg.PoolMaxConnLifetime > 0 {
		q.Set("pool_max_conn_lifetime", cfg.PoolMaxConnLifetime.String())
	}
	if cfg.PoolMaxConnIdleTime > 0 {
		q.Set("pool_max_conn_idle_time", cfg.PoolMaxConnIdleTime.String())
	}
	if cfg.PoolHealthCheckPeriod > 0 {
		q.Set("pool_health_check_period", cfg.PoolHealthCheckPeriod.String())
	}
	if cfg.PoolMaxConnLifetimeJitter > 0 {
		q.Set("pool_max_conn_lifetime_jitter", cfg.PoolMaxConnLifetimeJitter.String())
	}

	return &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
}

// Validate performs syntactic and semantic validation of Config. It does not
// verify network reachability or credentials. Checks required fields, ensures
// pool parameters are consistent, rejects negative durations, and validates
// host formatting (including IPv6 brackets if a port is present).
func (cfg Config) Validate() error {
	if cfg.User == "" {
		return errors.New("user is required")
	}
	if cfg.Password == "" {
		return errors.New("password is required")
	}
	if cfg.Host == "" {
		return errors.New("host is required")
	}
	if cfg.Name == "" {
		return errors.New("database name is required")
	}

	if _, err := url.Parse("postgres://" + cfg.Host); err != nil {
		return fmt.Errorf("invalid host: %w", err)
	}

	if cfg.PoolMinConns < 0 {
		return errors.New("pool_min_conns cannot be negative")
	}
	if cfg.PoolMaxConns < 0 {
		return errors.New("pool_max_conns cannot be negative")
	}
	if cfg.PoolMaxConns > 0 && cfg.PoolMinConns > cfg.PoolMaxConns {
		return errors.New("pool_min_conns cannot exceed pool_max_conns")
	}
	if cfg.PoolMaxConnLifetime < 0 {
		return errors.New("pool_max_conn_lifetime cannot be negative")
	}
	if cfg.PoolMaxConnIdleTime < 0 {
		return errors.New("pool_max_conn_idle_time cannot be negative")
	}
	if cfg.PoolHealthCheckPeriod < 0 {
		return errors.New("pool_health_check_period cannot be negative")
	}
	if cfg.PoolMaxConnLifetimeJitter < 0 {
		return errors.New("pool_max_conn_lifetime_jitter cannot be negative")
	}

	return nil
}
