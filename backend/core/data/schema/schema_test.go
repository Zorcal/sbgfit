package schema_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zorcal/sbgfit/backend/core/data/pgtest"
	"github.com/zorcal/sbgfit/backend/core/data/schema"
)

func TestSeed(t *testing.T) {
	ctx := context.Background()

	pool := pgtest.New(t, ctx)

	if err := schema.SeedData(ctx, pool); err != nil {
		t.Fatalf("first seed failed: %v", err)
	}

	initialRowCount := queryTotalRowCount(t, pool)
	if initialRowCount == 0 {
		t.Fatal("no data was seeded")
	}

	// Make sure seeding twice is a no-op.

	if err := schema.SeedData(ctx, pool); err != nil {
		t.Fatalf("second seed failed: %v", err)
	}

	finalRowCount := queryTotalRowCount(t, pool)
	if finalRowCount != initialRowCount {
		t.Errorf("total row count changed after second seed: got %d, want %d", finalRowCount, initialRowCount)
	}
}

// queryTotalRowCount returns the total number of rows across all tables in the
// sbgfit schema.
func queryTotalRowCount(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()

	tblNamesQ := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_type = 'BASE TABLE'
		AND table_schema = 'sbgfit'
		ORDER BY table_name;
	`
	rows, _ := pool.Query(t.Context(), tblNamesQ)
	tblNames, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		t.Fatalf("error collecting table names: %v", err)
	}

	var totalRows int
	for _, tblName := range tblNames {
		countQ := "SELECT COUNT(*) FROM sbgfit." + tblName
		row, _ := pool.Query(t.Context(), countQ)
		tableRows, err := pgx.CollectOneRow(row, pgx.RowTo[int])
		if err != nil {
			t.Fatalf("error counting rows in sbgfit.%s: %v", tblName, err)
		}
		totalRows += tableRows
	}

	return totalRows
}
