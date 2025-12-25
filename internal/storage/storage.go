package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func NewClient(ctx context.Context, connString string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{pool: pool}, nil
}

func (c *Client) Close() {
	c.pool.Close()
}
