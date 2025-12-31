package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool     *pgxpool.Pool
	ownsPool bool // true if we created the pool and should close it
}

// NewClient creates a new storage client with its own connection pool
func NewClient(ctx context.Context, connString string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{pool: pool, ownsPool: true}, nil
}

// NewClientFromPool creates a storage client using an existing connection pool
func NewClientFromPool(pool *pgxpool.Pool) *Client {
	return &Client{pool: pool, ownsPool: false}
}

// Close closes the connection pool only if we own it
func (c *Client) Close() {
	if c.ownsPool && c.pool != nil {
		c.pool.Close()
	}
}
