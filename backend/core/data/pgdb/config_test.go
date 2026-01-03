package pgdb_test

import (
	"testing"
	"time"

	"github.com/zorcal/sbgfit/backend/core/data/pgdb"
)

func TestConfig_ParsePgxConfig(t *testing.T) {
	tests := []struct {
		name        string
		in          pgdb.Config
		wantConnStr string
	}{
		{
			name: "minimal config",
			in: pgdb.Config{
				User:     "testuser",
				Password: "testpass",
				Host:     "localhost:5432",
				Name:     "testdb",
			},
			wantConnStr: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=require&timezone=utc",
		},
		{
			name: "all fields configured",
			in: pgdb.Config{
				User:                      "testuser",
				Password:                  "testpass",
				Host:                      "localhost:5432",
				Name:                      "testdb",
				Schema:                    "public",
				DisableTLS:                true,
				PoolMaxConns:              10,
				PoolMinConns:              2,
				PoolMaxConnLifetime:       time.Hour,
				PoolMaxConnIdleTime:       30 * time.Minute,
				PoolHealthCheckPeriod:     5 * time.Minute,
				PoolMaxConnLifetimeJitter: 10 * time.Second,
			},
			wantConnStr: "postgres://testuser:testpass@localhost:5432/testdb?pool_health_check_period=5m0s&pool_max_conn_idle_time=30m0s&pool_max_conn_lifetime=1h0m0s&pool_max_conn_lifetime_jitter=10s&pool_max_conns=10&pool_min_conns=2&search_path=public&sslmode=disable&timezone=utc",
		},
		{
			name: "IPv6 host",
			in: pgdb.Config{
				User:     "testuser",
				Password: "testpass",
				Host:     "[::1]:5432",
				Name:     "testdb",
			},
			wantConnStr: "postgres://testuser:testpass@[::1]:5432/testdb?sslmode=require&timezone=utc",
		},
		{
			name: "TLS enabled",
			in: pgdb.Config{
				User:       "testuser",
				Password:   "testpass",
				Host:       "localhost:5432",
				Name:       "testdb",
				DisableTLS: false,
			},
			wantConnStr: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=require&timezone=utc",
		},
		{
			name: "empty schema",
			in: pgdb.Config{
				User:     "testuser",
				Password: "testpass",
				Host:     "localhost:5432",
				Name:     "testdb",
				Schema:   "",
			},
			wantConnStr: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=require&timezone=utc",
		},
		{
			name: "non-empty schema",
			in: pgdb.Config{
				User:     "testuser",
				Password: "testpass",
				Host:     "localhost:5432",
				Name:     "testdb",
				Schema:   "public",
			},
			wantConnStr: "postgres://testuser:testpass@localhost:5432/testdb?search_path=public&sslmode=require&timezone=utc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.ParsePgxConfig()
			if err != nil {
				t.Fatalf("Config.ParsePgxConfig(%+v) error = %v, want nil", tt.in, err)
			}
			if got == nil {
				t.Fatalf("Config.ParsePgxConfig(%+v) = nil, want *pgxpool.Config", tt.in)
			}

			if got.ConnConfig == nil {
				t.Fatalf("Config.ParsePgxConfig(%+v) returned config with nil ConnConfig", tt.in)
			}

			gotConnStr := got.ConnString()
			if gotConnStr != tt.wantConnStr {
				t.Errorf("Config.ParsePgxConfig(%+v) ConnString() = %q, want %q", tt.in, gotConnStr, tt.wantConnStr)
			}
		})
	}
}

func TestConfig_ParsePgxConfig_error(t *testing.T) {
	tests := []struct {
		name string
		in   pgdb.Config
	}{
		{
			name: "missing user",
			in: pgdb.Config{
				Password: "testpass",
				Host:     "localhost:5432",
				Name:     "testdb",
			},
		},
		{
			name: "missing password",
			in: pgdb.Config{
				User: "testuser",
				Host: "localhost:5432",
				Name: "testdb",
			},
		},
		{
			name: "missing host",
			in: pgdb.Config{
				User:     "testuser",
				Password: "testpass",
				Name:     "testdb",
			},
		},
		{
			name: "missing database name",
			in: pgdb.Config{
				User:     "testuser",
				Password: "testpass",
				Host:     "localhost:5432",
			},
		},
		{
			name: "invalid host format",
			in: pgdb.Config{
				User:     "testuser",
				Password: "testpass",
				Host:     "invalid host with spaces",
				Name:     "testdb",
			},
		},
		{
			name: "negative pool max conns",
			in: pgdb.Config{
				User:         "testuser",
				Password:     "testpass",
				Host:         "localhost:5432",
				Name:         "testdb",
				PoolMaxConns: -1,
			},
		},
		{
			name: "negative pool min conns",
			in: pgdb.Config{
				User:         "testuser",
				Password:     "testpass",
				Host:         "localhost:5432",
				Name:         "testdb",
				PoolMinConns: -1,
			},
		},
		{
			name: "min conns greater than max conns",
			in: pgdb.Config{
				User:         "testuser",
				Password:     "testpass",
				Host:         "localhost:5432",
				Name:         "testdb",
				PoolMaxConns: 2,
				PoolMinConns: 5,
			},
		},
		{
			name: "negative max conn lifetime",
			in: pgdb.Config{
				User:                "testuser",
				Password:            "testpass",
				Host:                "localhost:5432",
				Name:                "testdb",
				PoolMaxConnLifetime: -time.Hour,
			},
		},
		{
			name: "negative max conn idle time",
			in: pgdb.Config{
				User:                "testuser",
				Password:            "testpass",
				Host:                "localhost:5432",
				Name:                "testdb",
				PoolMaxConnIdleTime: -time.Minute,
			},
		},
		{
			name: "negative health check period",
			in: pgdb.Config{
				User:                  "testuser",
				Password:              "testpass",
				Host:                  "localhost:5432",
				Name:                  "testdb",
				PoolHealthCheckPeriod: -time.Second,
			},
		},
		{
			name: "negative lifetime jitter",
			in: pgdb.Config{
				User:                      "testuser",
				Password:                  "testpass",
				Host:                      "localhost:5432",
				Name:                      "testdb",
				PoolMaxConnLifetimeJitter: -time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.in.ParsePgxConfig(); err == nil {
				t.Errorf("Config.ParsePgxConfig(%+v) error = nil, want error", tt.in)
			}
		})
	}
}
