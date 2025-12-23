package chunker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChunkDocument(t *testing.T) {
	// Read a sample MDX file
	testFile := filepath.Join("..", "..", "docs", "strudel", "learn", "notes.mdx")
	content, err := os.ReadFile(testFile)

	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	opts := DefaultOptions()

	chunks, err := ChunkDocument(string(content), "notes.mdx", opts)
	if err != nil {
		t.Fatalf("ChunkDocument failed: %v", err)
	}

	if len(chunks) == 0 {
		t.Fatal("Expected at least one chunk, got 0")
	}

	t.Logf("Generated %d chunks from notes.mdx\n", len(chunks))
	t.Logf("========================================\n")

	for i, chunk := range chunks {
		t.Logf("\n--- Chunk %d ---", i+1)
		t.Logf("Page: %s", chunk.PageName)
		t.Logf("URL: %s", chunk.PageURL)
		t.Logf("Section: %s", chunk.SectionTitle)
		t.Logf("Content length: %d chars (~%d tokens)", len(chunk.Content), estimateTokens(chunk.Content))
		t.Logf("Metadata: %+v", chunk.Metadata)
		t.Logf("\nContent preview (first 200 chars):")
		preview := chunk.Content

		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}

		t.Logf("%s\n", preview)
	}
}
