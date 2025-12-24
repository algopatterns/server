package chunker

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DefaultOptions() ChunkOptions {
	return ChunkOptions{
		MaxTokens:       800,
		OverlapTokens:   100,
		PreserveHeaders: true,
	}
}

func ChunkDocument(content, pageName string, opts ChunkOptions) ([]Chunk, error) {
	metadata := extractFrontmatter(content)
	content = frontmatterRegex.ReplaceAllString(content, "")
	content = importRegex.ReplaceAllString(content, "")
	content = stripMDXComponents(content)

	pageURL := generateURL(pageName)
	sections := splitByHeaders(content)

	var chunks []Chunk

	for _, section := range sections {
		if estimateTokens(section.Content) <= opts.MaxTokens {
			chunks = append(chunks, Chunk{
				PageName:     pageName,
				PageURL:      pageURL,
				SectionTitle: section.Title,
				Content:      strings.TrimSpace(section.Content),
				Metadata:     metadata,
			})

			continue
		}

		subChunks := splitLargeSection(section, opts)

		for _, subChunk := range subChunks {
			chunks = append(chunks, Chunk{
				PageName:     pageName,
				PageURL:      pageURL,
				SectionTitle: section.Title,
				Content:      strings.TrimSpace(subChunk),
				Metadata:     metadata,
			})
		}
	}

	return chunks, nil
}

// ChunkDocuments discovers all markdown files in a directory and chunks them
// Returns chunks and a slice of errors encountered (one per failed file)
func ChunkDocuments(docsPath string) ([]Chunk, []error) {
	opts := DefaultOptions()
	var allChunks []Chunk
	var errors []error
	fileCount := 0

	// walk the directory tree to find all markdown files
	walkErr := filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("warning: error accessing path %s: %v", path, err)
			errors = append(errors, fmt.Errorf("path %s: %w", path, err))
			return nil // continue walking
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".md" && ext != ".mdx" {
			return nil
		}

		fileCount++

		// read file content
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("warning: failed to read file %s: %v", path, err)
			errors = append(errors, fmt.Errorf("read %s: %w", path, err))
			return nil // continue with other files
		}

		// get relative page name
		pageName, err := filepath.Rel(docsPath, path)
		if err != nil {
			pageName = filepath.Base(path)
		}

		// chunk the document
		chunks, err := ChunkDocument(string(content), pageName, opts)
		if err != nil {
			log.Printf("warning: failed to chunk document %s: %v", path, err)
			errors = append(errors, fmt.Errorf("chunk %s: %w", path, err))
			return nil // continue with other files
		}

		allChunks = append(allChunks, chunks...)

		return nil
	})

	if walkErr != nil {
		errors = append(errors, fmt.Errorf("walk error: %w", walkErr))
	}

	log.Printf("processed %d markdown files, generated %d chunks, encountered %d errors",
		fileCount, len(allChunks), len(errors))

	return allChunks, errors
}
