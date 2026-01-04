package pgtest

import (
	"context"
	"fmt"
	"os"
	"testing"
)

const pgTemplateName = "pgtest_template"

func TestMain(m *testing.M) {
	ctx := context.Background()

	cleanup, err := NewTemplate(ctx, pgTemplateName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "TestMain: create new posgres template: %s\n", err)
		os.Exit(1)
	}

	code := m.Run()

	if err := cleanup(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "TestMain: cleanup: %s\n", err)
		os.Exit(1)
	}

	os.Exit(code)
}

func TestTemplate(t *testing.T) {
	ctx := context.Background()

	dbName := t.Name()
	pool := FromTemplate(t, ctx, pgTemplateName, dbName)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Ping error: %s", err)
	}
}
