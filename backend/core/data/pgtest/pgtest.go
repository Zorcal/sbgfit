// Package pgtest provides PostgreSQL test helpers.
//
// It supports:
//   - Creating temporary test databases from a template (`FromTemplate`)
//   - Creating reusable database templates for package-level tests (`NewTemplate`)
//
// All database operations assume PostgreSQL is running via docker-compose.
package pgtest

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zorcal/sbgfit/backend/core/data/pgdb"
	"github.com/zorcal/sbgfit/backend/core/data/schema"
)

// Static database configuration.
const (
	hostPort   = "sbgfit-postgres:5432" // See docker-compose.yml at repository root.
	username   = "postgres"
	password   = "postgres"
	disableTLS = true
)

// FromTemplate creates a temporary test database from an existing template.
func FromTemplate(t *testing.T, ctx context.Context, tmplName, dbName string) *pgxpool.Pool {
	t.Helper()

	if err := validateDBName(dbName); err != nil {
		t.Fatalf("validate database name: %s", err)
	}

	mcfg, err := dbCfg("postgres").ParsePgxConfig()
	if err != nil {
		t.Fatalf("parse manager database config: %v", err)
	}

	mpool, err := pgdb.NewPool(ctx, mcfg)
	if err != nil {
		t.Fatalf("create database manager pool: %s", err)
	}
	t.Cleanup(func() { mpool.Close() })

	if err := pgdb.StatusCheck(ctx, mpool); err != nil {
		t.Fatalf(`status check database manager: %s

		Did you remember to run 'docker-compose up -d' at the repository root?
		`, err)
	}

	if _, err := mpool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %q", dbName)); err != nil {
		t.Fatalf("drop database: %s", err)
	}
	if _, err := mpool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %q TEMPLATE %q", dbName, tmplName)); err != nil {
		t.Fatalf("create database %s from template %s: %s", dbName, tmplName, err)
	}
	t.Cleanup(func() {
		if _, err := mpool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %q", dbName)); err != nil {
			t.Fatalf("cleanup: drop database: %s", err)
		}
	})

	cfg, err := dbCfg(dbName).ParsePgxConfig()
	if err != nil {
		t.Fatalf("parse db pgx config: %v", err)
	}

	pool, err := pgdb.NewPool(ctx, cfg)
	if err != nil {
		t.Fatalf("create database pool %s: %s", dbName, err)
	}
	t.Cleanup(func() { pool.Close() })

	if err := pgdb.StatusCheck(ctx, pool); err != nil {
		t.Fatalf("status check database: %s", err)
	}

	return pool
}

// NewTemplate creates a reusable database template for tests. Intended for
// TestMain usage.
func NewTemplate(ctx context.Context, dbName string) (cleanup func(context.Context) error, retErr error) {
	if err := validateDBName(dbName); err != nil {
		return nil, fmt.Errorf("validate database name: %w", err)
	}

	mcfg, err := dbCfg("postgres").ParsePgxConfig()
	if err != nil {
		return nil, fmt.Errorf("parse manager database config: %w", err)
	}

	mpool, err := pgdb.NewPool(ctx, mcfg)
	if err != nil {
		return nil, fmt.Errorf("create database manager pool: %w", err)
	}
	defer func() {
		if retErr != nil && mpool != nil {
			mpool.Close()
		}
	}()

	if err := pgdb.StatusCheck(ctx, mpool); err != nil {
		return nil, fmt.Errorf("status check database manager: %w", err)
	}

	var exists bool
	if err := mpool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)`, dbName).Scan(&exists); err != nil {
		return nil, fmt.Errorf("existence check: %w", err)
	}
	if exists {
		if _, err := mpool.Exec(ctx, fmt.Sprintf("ALTER DATABASE %q WITH IS_TEMPLATE=false", dbName)); err != nil {
			return nil, fmt.Errorf("untemplate database: %w", err)
		}
		if _, err := mpool.Exec(ctx, fmt.Sprintf("DROP DATABASE %q", dbName)); err != nil {
			return nil, fmt.Errorf("drop database: %w", err)
		}
	}

	if _, err := mpool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %q", dbName)); err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}

	if err := schema.Migrate(ctx, dbCfg(dbName)); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	if _, err := mpool.Exec(ctx, fmt.Sprintf("ALTER DATABASE %q WITH IS_TEMPLATE=true", dbName)); err != nil {
		return nil, fmt.Errorf("make template: %w", err)
	}

	cleanup = func(ctx context.Context) error {
		var err error
		if _, execErr := mpool.Exec(ctx, fmt.Sprintf("ALTER DATABASE %q WITH IS_TEMPLATE=false", dbName)); execErr != nil {
			err = errors.Join(err, execErr)
		}
		if _, execErr := mpool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %q", dbName)); execErr != nil {
			err = errors.Join(err, execErr)
		}
		mpool.Close()
		return err
	}

	return cleanup, nil
}

func dbCfg(dbName string) pgdb.Config {
	return pgdb.Config{
		User:       cmp.Or(os.Getenv("POSTGRES_USER"), username),
		Password:   cmp.Or(os.Getenv("POSTGRES_PASSWORD"), password),
		Host:       cmp.Or(os.Getenv("POSTGRES_HOST"), hostPort),
		Name:       dbName,
		DisableTLS: disableTLS,
	}
}

func validateDBName(name string) error {
	if name == "" || len(name) > 63 {
		return errors.New("database name must be between 1 and 63 characters long")
	}
	if !unicode.IsLetter(rune(name[0])) {
		return errors.New("database name must start with a letter")
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return errors.New("database name can only contain letters, digits, and underscores")
		}
	}
	return nil
}
