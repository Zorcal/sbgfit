// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	"embed"
	"fmt"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

//go:embed seed.sql
var seedSQL string

// Migrate attempts to bring the database up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, connStr string) error {
	connURL, err := url.Parse(connStr)
	if err != nil {
		return fmt.Errorf("parse conn URL: %w", err)
	}

	db := dbmate.New(connURL)
	db.FS = migrationsFS
	db.MigrationsDir = []string{"./migrations"}

	if err := db.CreateAndMigrate(); err != nil {
		return fmt.Errorf("create and migrate: %w", err)
	}

	return nil
}

// SeedData seeds the database with static seed data.
func SeedData(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, seedSQL); err != nil {
		return fmt.Errorf("exec seed SQL: %w", err)
	}
	return nil
}
