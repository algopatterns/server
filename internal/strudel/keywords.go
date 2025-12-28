package strudel

import "strings"

// returns sensible defaults for keyword extraction
func DefaultKeywordOptions() KeywordOptions {
	return KeywordOptions{
		MaxKeywords:      10,
		IncludeSounds:    true,
		IncludeNotes:     true,
		IncludeFunctions: true,
		IncludeScales:    true,
		Deduplicate:      true,
	}
}

func ExtractKeywords(code string) string {
	return ExtractKeywordsWithOptions(code, DefaultKeywordOptions())
}

func ExtractKeywordsWithOptions(code string, opts KeywordOptions) string {
	if code == "" {
		return ""
	}

	keywords := []string{}

	parsed := Parse(code)

	if opts.IncludeSounds {
		keywords = append(keywords, parsed.Sounds...)
	}

	if opts.IncludeNotes {
		keywords = append(keywords, parsed.Notes...)
	}

	if opts.IncludeFunctions {
		keywords = append(keywords, parsed.Functions...)
	}

	if opts.IncludeScales {
		keywords = append(keywords, parsed.Scales...)
	}

	if opts.Deduplicate {
		keywords = UniqueStrings(keywords)
	}

	if opts.MaxKeywords > 0 && len(keywords) > opts.MaxKeywords {
		keywords = keywords[:opts.MaxKeywords]
	}

	return strings.Join(keywords, " ")
}

// removes duplicates while preserving order
func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, s := range slice {
		if !seen[s] {
			result = append(result, s)
			seen[s] = true
		}
	}

	return result
}
