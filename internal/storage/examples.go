package storage

import (
	"context"
	"fmt"

	"github.com/algorave/server/internal/examples"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

// ClearAllExamples deletes all existing examples from the database
func (c *Client) ClearAllExamples(ctx context.Context) error {
	_, err := c.pool.Exec(ctx, "DELETE FROM example_strudels")
	if err != nil {
		return fmt.Errorf("failed to clear examples: %w", err)
	}

	return nil
}

// InsertExamplesBatch inserts multiple examples in a single transaction
func (c *Client) InsertExamplesBatch(ctx context.Context, examples []examples.Example, embeddings [][]float32) error {
	if len(examples) != len(embeddings) {
		return fmt.Errorf("examples and embeddings length mismatch")
	}

	if len(examples) == 0 {
		return nil
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	query := `
		INSERT INTO example_strudels (title, description, code, tags, embedding, url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for index, example := range examples {
		batch.Queue(query,
			example.Title,
			example.Description,
			example.Code,
			example.Tags,
			pgvector.NewVector(embeddings[index]),
			example.SourceURL,
		)
	}

	br := tx.SendBatch(ctx, batch)

	for i := range examples {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("failed to insert example %d: %w", i, err)
		}
	}

	// must close batch results before committing
	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetExampleCount returns the total number of examples in the database
func (c *Client) GetExampleCount(ctx context.Context) (int, error) {
	var count int

	err := c.pool.QueryRow(ctx, "SELECT COUNT(*) FROM example_strudels").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get example count: %w", err)
	}

	return count, nil
}
