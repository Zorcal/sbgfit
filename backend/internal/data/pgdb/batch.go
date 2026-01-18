package pgdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zorcal/sbgfit/backend/internal/telemetry"
)

// Batch wraps a pgx.Batch together with the pgxpool.Pool it will be executed on.
//
// It is used to queue multiple database operations and execute them together
// as a single batch.
type Batch struct {
	b *pgx.Batch
	p *pgxpool.Pool
}

func newBatch(p *pgxpool.Pool) *Batch {
	return &Batch{
		b: &pgx.Batch{},
		p: p,
	}
}

// RunBatch creates a new Batch, passes it to f for query queueing, and then
// executes the batch against the provided pool.
//
// If f returns an error, the batch is not sent. If sending or closing the
// batch results fails, RunBatch returns an error.
func RunBatch(ctx context.Context, p *pgxpool.Pool, queueFunc func(ctx context.Context, b *Batch) error) error {
	ctx, span := telemetry.StartSpan(ctx, "pgdb.BatchTx")
	defer span.End()

	b := newBatch(p)

	if err := queueFunc(ctx, b); err != nil {
		return fmt.Errorf("queueFunc: %w", err)
	}

	span.AddEvent("Batch started")

	result := b.p.SendBatch(ctx, b.b)
	if err := result.Close(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("close batch result: %w", err)
	}

	span.AddEvent("Batch successfully finished")

	return nil
}
