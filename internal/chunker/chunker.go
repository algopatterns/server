package chunker

import (
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
