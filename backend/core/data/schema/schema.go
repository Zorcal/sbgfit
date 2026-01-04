// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	"embed"
	"fmt"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

	"github.com/zorcal/sbgfit/backend/core/data/pgdb"
)

//go:embed migrations/*.sql
var fs embed.FS

// Migrate attempts to bring the database up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, cfg pgdb.Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	db := dbmate.New(cfg.URL())
	db.FS = fs
	db.MigrationsDir = []string{"./migrations"}

	if err := db.CreateAndMigrate(); err != nil {
		return fmt.Errorf("create and migrate: %w", err)
	}

	return nil
}
