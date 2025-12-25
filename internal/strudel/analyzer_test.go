package strudel

import (
	"strings"
	"testing"
)

func TestAnalyzeCode_Drums(t *testing.T) {
	code := `sound("bd hh sd").fast(4)`

	analysis := AnalyzeCode(code)

	if !contains(analysis.SoundTags, "drums") {
		t.Errorf("Expected 'drums' in SoundTags, got: %v", analysis.SoundTags)
	}

	if !contains(analysis.SoundTags, "percussion") {
		t.Errorf("Expected 'percussion' in SoundTags, got: %v", analysis.SoundTags)
	}

	if analysis.Complexity < 1 {
		t.Errorf("Expected Complexity > 0, got: %d", analysis.Complexity)
	}
}

func TestAnalyzeCode_Melody(t *testing.T) {
	code := `note("c e g").scale("minor")`

	analysis := AnalyzeCode(code)

	if !contains(analysis.MusicalTags, "melody") {
		t.Errorf("Expected 'melody' in MusicalTags, got: %v", analysis.MusicalTags)
	}

	if !contains(analysis.MusicalTags, "melodic") {
		t.Errorf("Expected 'melodic' in MusicalTags, got: %v", analysis.MusicalTags)
	}

	if !contains(analysis.MusicalTags, "scales") {
		t.Errorf("Expected 'scales' in MusicalTags, got: %v", analysis.MusicalTags)
	}
}

func TestAnalyzeCode_Effects(t *testing.T) {
	code := `sound("bd").delay(0.25).room(0.5).lpf(1000)`

	analysis := AnalyzeCode(code)

	if !contains(analysis.EffectTags, "delay") {
		t.Errorf("Expected 'delay' in EffectTags, got: %v", analysis.EffectTags)
	}

	if !contains(analysis.EffectTags, "reverb") {
		t.Errorf("Expected 'reverb' in EffectTags, got: %v", analysis.EffectTags)
	}

	if !contains(analysis.EffectTags, "filter") {
		t.Errorf("Expected 'filter' in EffectTags, got: %v", analysis.EffectTags)
	}
}

func TestAnalyzeCode_Complexity(t *testing.T) {
	tests := []struct {
		name             string
		code             string
		minComplexity    int
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:             "simple pattern",
			code:             `sound("bd")`,
			minComplexity:    1,
			shouldContain:    []string{"simple", "beginner-friendly"},
			shouldNotContain: []string{"complex", "layered", "advanced"},
		},
		{
			name: "complex layered pattern",
			code: `
			let pat1 = sound("bd").fast(4)
			let pat2 = sound("hh").fast(8)
			let pat3 = note("c e g")
			pat1.stack(pat2).stack(pat3).stack(sound("sd"))
			`,
			minComplexity:    4,
			shouldContain:    []string{"layered"},
			shouldNotContain: []string{"simple", "beginner-friendly"},
		},
		{
			name:             "arranged pattern",
			code:             `arrange(1, [sound("bd"), sound("hh")]).stack(note("c"))`,
			minComplexity:    2,
			shouldContain:    []string{"arranged", "structured"},
			shouldNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeCode(tt.code)

			if analysis.Complexity < tt.minComplexity {
				t.Errorf("Expected Complexity >= %d, got: %d", tt.minComplexity, analysis.Complexity)
			}

			for _, tag := range tt.shouldContain {
				if !contains(analysis.ComplexityTags, tag) {
					t.Errorf("Expected '%s' in ComplexityTags, got: %v", tag, analysis.ComplexityTags)
				}
			}

			for _, tag := range tt.shouldNotContain {
				if contains(analysis.ComplexityTags, tag) {
					t.Errorf("Did not expect '%s' in ComplexityTags, got: %v", tag, analysis.ComplexityTags)
				}
			}
		})
	}
}

func TestGenerateTags(t *testing.T) {
	code := `sound("bd").fast(4).stack(note("c e g"))`

	analysis := AnalyzeCode(code)

	tests := []struct {
		name          string
		category      string
		existingTags  []string
		shouldContain []string
	}{
		{
			name:          "with category",
			category:      "techno",
			existingTags:  []string{},
			shouldContain: []string{"techno", "drums", "percussion"},
		},
		{
			name:          "with existing tags",
			category:      "",
			existingTags:  []string{"tutorial", "beginner"},
			shouldContain: []string{"tutorial", "beginner"},
		},
		{
			name:          "with category and existing tags",
			category:      "ambient",
			existingTags:  []string{"experimental"},
			shouldContain: []string{"ambient", "experimental"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := GenerateTags(analysis, tt.category, tt.existingTags)

			for _, tag := range tt.shouldContain {
				if !contains(tags, tag) && !contains(tags, strings.ToLower(tag)) {
					t.Errorf("Expected '%s' in tags, got: %v", tag, tags)
				}
			}
		})
	}
}

func TestCalculateComplexity(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "minimal code",
			code:     `sound("bd")`,
			expected: 1,
		},
		{
			name:     "medium code with stack",
			code:     `sound("bd").stack(sound("hh"))`,
			expected: 3,
		},
		{
			name: "complex code",
			code: `
			let a = 1
			let b = 2
			let c = 3
			let d = 4
			let e = 5
			let f = 6
			sound("bd").stack(sound("hh")).stack(sound("sd")).stack(sound("cp"))
			`,
			expected: 5, // base 2 (length) + stack 3 + var 2 - limited by scoring
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := Parse(tt.code)
			complexity := calculateComplexity(tt.code, parsed)

			if complexity != tt.expected {
				t.Errorf("calculateComplexity() = %d, want %d", complexity, tt.expected)
			}
		})
	}
}
