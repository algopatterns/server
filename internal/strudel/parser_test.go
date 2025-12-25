package strudel

import (
	"reflect"
	"testing"
)

func TestExtractSounds(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "single sound",
			code:     `sound("bd")`,
			expected: []string{"bd"},
		},
		{
			name:     "multiple sounds in one call",
			code:     `sound("bd hh sd")`,
			expected: []string{"bd", "hh", "sd"},
		},
		{
			name:     "multiple sound calls",
			code:     `sound("bd").stack(sound("hh"))`,
			expected: []string{"bd", "hh"},
		},
		{
			name:     "short form s()",
			code:     `s("bd").stack(s("hh"))`,
			expected: []string{"bd", "hh"},
		},
		{
			name:     "no sounds",
			code:     `note("c e g")`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSounds(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractSounds() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractNotes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "single note",
			code:     `note("c")`,
			expected: []string{"c"},
		},
		{
			name:     "multiple notes",
			code:     `note("c e g")`,
			expected: []string{"c", "e", "g"},
		},
		{
			name:     "multiple note calls",
			code:     `note("c e").stack(note("g a"))`,
			expected: []string{"c", "e", "g", "a"},
		},
		{
			name:     "no notes",
			code:     `sound("bd")`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractNotes(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractNotes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractFunctions(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "single function",
			code:     `sound("bd").fast(2)`,
			expected: []string{"fast"},
		},
		{
			name:     "multiple functions",
			code:     `sound("bd").fast(2).slow(4).stack()`,
			expected: []string{"fast", "slow", "stack"},
		},
		{
			name:     "no functions",
			code:     `sound("bd")`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractFunctions(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractFunctions() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "let declaration",
			code:     `let pat1 = sound("bd")`,
			expected: []string{"pat1"},
		},
		{
			name:     "const declaration",
			code:     `const rhythm = sound("hh")`,
			expected: []string{"rhythm"},
		},
		{
			name:     "multiple declarations",
			code:     `let x = 1; const y = 2; var z = 3`,
			expected: []string{"x", "y", "z"},
		},
		{
			name:     "no variables",
			code:     `sound("bd").fast(2)`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractVariables(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractVariables() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountPatterns(t *testing.T) {
	code := `sound("bd").stack(sound("hh")).every(4).stack(note("c"))`

	patterns := CountPatterns(code)

	if patterns["stack"] != 2 {
		t.Errorf("CountPatterns() stack = %d, want 2", patterns["stack"])
	}

	if patterns["every"] != 1 {
		t.Errorf("CountPatterns() every = %d, want 1", patterns["every"])
	}

	if patterns["arrange"] != 0 {
		t.Errorf("CountPatterns() arrange = %d, want 0", patterns["arrange"])
	}
}

func TestParse(t *testing.T) {
	code := `
let pat1 = sound("bd hh sd").fast(2)
note("c e g").scale("minor").stack(pat1)
	`

	parsed := Parse(code)

	// Check sounds
	expectedSounds := []string{"bd", "hh", "sd"}
	if !reflect.DeepEqual(parsed.Sounds, expectedSounds) {
		t.Errorf("Parse().Sounds = %v, want %v", parsed.Sounds, expectedSounds)
	}

	// Check notes
	expectedNotes := []string{"c", "e", "g"}
	if !reflect.DeepEqual(parsed.Notes, expectedNotes) {
		t.Errorf("Parse().Notes = %v, want %v", parsed.Notes, expectedNotes)
	}

	// Check functions
	if len(parsed.Functions) < 2 {
		t.Errorf("Parse().Functions has %d items, want at least 2", len(parsed.Functions))
	}

	// Check variables
	if len(parsed.Variables) != 1 || parsed.Variables[0] != "pat1" {
		t.Errorf("Parse().Variables = %v, want [pat1]", parsed.Variables)
	}

	// Check patterns
	if parsed.Patterns["stack"] != 1 {
		t.Errorf("Parse().Patterns[stack] = %d, want 1", parsed.Patterns["stack"])
	}
}
