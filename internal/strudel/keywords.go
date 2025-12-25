package strudel

import "strings"

// configures keyword extraction behavior
type KeywordOptions struct {
	MaxKeywords      int  // limit total keywords (default: 10)
	IncludeSounds    bool // include sound names (default: true)
	IncludeNotes     bool // include note names (default: true)
	IncludeFunctions bool // include function names (default: true)
	IncludeScales    bool // include scale names (default: true)
	Deduplicate      bool // remove duplicates (default: true)
}

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

// extracts keywords with default options
// returns a space-separated keyword string suitable for search queries
func ExtractKeywords(code string) string {
	return ExtractKeywordsWithOptions(code, DefaultKeywordOptions())
}

// extracts keywords with custom options
func ExtractKeywordsWithOptions(code string, opts KeywordOptions) string {
	if code == "" {
		return ""
	}

	keywords := []string{}

	// parse code once
	parsed := Parse(code)

	// collect keywords based on options
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

	// deduplicate if requested
	if opts.Deduplicate {
		keywords = UniqueStrings(keywords)
	}

	// limit to max keywords
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
