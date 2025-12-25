package retriever

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Embedder generates embeddings from text (narrow interface)
type Embedder interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}

// QueryTransformer expands user queries with technical keywords (narrow interface)
type QueryTransformer interface {
	TransformQuery(ctx context.Context, userQuery string) (string, error)
}

// Client performs vector similarity search on documentation and examples
type Client struct {
	db          *pgxpool.Pool
	embedder    Embedder
	transformer QueryTransformer
	topK        int
}

// SearchResult represents a document chunk from vector search
type SearchResult struct {
	ID           string
	PageName     string
	PageURL      string
	SectionTitle string
	Content      string
	Similarity   float32
	Metadata     map[string]interface{}
}

// ExampleResult represents an example Strudel from vector search
type ExampleResult struct {
	ID          string
	Title       string
	Description string
	Code        string
	Tags        []string
	URL         string
	Similarity  float32
}
