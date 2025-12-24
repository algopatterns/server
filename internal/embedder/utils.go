package embedder

import (
	"path/filepath"
)

const (
	openaiEmbeddingsURL = "https://api.openai.com/v1/embeddings"
	embeddingModel      = "text-embedding-3-small"
	embeddingDimensions = 1536
)

// extracts a relative page name from the full file path
func getPageName(fullPath, rootPath string) string {
	relPath, err := filepath.Rel(rootPath, fullPath)
	if err != nil {
		return filepath.Base(fullPath)
	}

	return relPath
}
