package strudel

import (
	"strings"
	"testing"
)

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		contains []string
	}{
		{
			name:     "basic pattern",
			code:     `sound("bd").fast(2)`,
			contains: []string{"bd", "fast"},
		},
		{
			name:     "complex pattern",
			code:     `sound("bd hh sd").fast(2).stack(note("c e g"))`,
			contains: []string{"bd", "hh", "sd", "fast", "stack", "c", "e", "g"},
		},
		{
			name:     "with scales",
			code:     `note("c e g").scale("minor")`,
			contains: []string{"c", "e", "g", "scale", "minor"},
		},
		{
			name:     "empty code",
			code:     "",
			contains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractKeywords(tt.code)

			if tt.code == "" && result != "" {
				t.Errorf("ExtractKeywords(%q) = %q, want empty string", tt.code, result)
				return
			}

			for _, keyword := range tt.contains {
				if !strings.Contains(result, keyword) {
					t.Errorf("ExtractKeywords() result %q does not contain %q", result, keyword)
				}
			}
		})
	}
}

func TestExtractKeywordsWithOptions(t *testing.T) {
	code := `sound("bd").fast(2).stack(note("c e g")).scale("minor")`

	t.Run("only sounds", func(t *testing.T) {
		opts := KeywordOptions{
			MaxKeywords:      10,
			IncludeSounds:    true,
			IncludeNotes:     false,
			IncludeFunctions: false,
			IncludeScales:    false,
			Deduplicate:      true,
		}

		result := ExtractKeywordsWithOptions(code, opts)

		if !strings.Contains(result, "bd") {
			t.Errorf("Result should contain 'bd', got: %s", result)
		}

		if strings.Contains(result, "fast") {
			t.Errorf("Result should not contain 'fast', got: %s", result)
		}
	})

	t.Run("max keywords limit", func(t *testing.T) {
		opts := KeywordOptions{
			MaxKeywords:      3,
			IncludeSounds:    true,
			IncludeNotes:     true,
			IncludeFunctions: true,
			IncludeScales:    true,
			Deduplicate:      true,
		}

		result := ExtractKeywordsWithOptions(code, opts)
		keywords := strings.Fields(result)

		if len(keywords) > 3 {
			t.Errorf("Expected max 3 keywords, got %d: %v", len(keywords), keywords)
		}
	})

	t.Run("deduplication", func(t *testing.T) {
		duplicateCode := `sound("bd bd bd").fast(2).fast(3)`

		opts := KeywordOptions{
			MaxKeywords:      10,
			IncludeSounds:    true,
			IncludeNotes:     false,
			IncludeFunctions: true,
			IncludeScales:    false,
			Deduplicate:      true,
		}

		result := ExtractKeywordsWithOptions(duplicateCode, opts)
		keywords := strings.Fields(result)

		// Count occurrences of "bd"
		bdCount := 0
		for _, kw := range keywords {
			if kw == "bd" {
				bdCount++
			}
		}

		if bdCount > 1 {
			t.Errorf("Expected 'bd' to appear once, got %d times in: %v", bdCount, keywords)
		}
	})
}

func TestDefaultKeywordOptions(t *testing.T) {
	opts := DefaultKeywordOptions()

	if opts.MaxKeywords != 10 {
		t.Errorf("DefaultKeywordOptions().MaxKeywords = %d, want 10", opts.MaxKeywords)
	}

	if !opts.IncludeSounds {
		t.Error("DefaultKeywordOptions().IncludeSounds should be true")
	}

	if !opts.IncludeNotes {
		t.Error("DefaultKeywordOptions().IncludeNotes should be true")
	}

	if !opts.IncludeFunctions {
		t.Error("DefaultKeywordOptions().IncludeFunctions should be true")
	}

	if !opts.Deduplicate {
		t.Error("DefaultKeywordOptions().Deduplicate should be true")
	}
}
