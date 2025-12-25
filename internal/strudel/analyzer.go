package strudel

import (
	"regexp"
	"slices"
	"strings"
)

// contains semantic analysis results
type CodeAnalysis struct {
	// categorized tags
	SoundTags      []string // ["drums", "synth", "bass"]
	EffectTags     []string // ["delay", "reverb", "filter"]
	MusicalTags    []string // ["melody", "chords", "rhythm"]
	ComplexityTags []string // ["layered", "advanced", "simple"]

	// metrics
	Complexity    int // 0-10 score
	LineCount     int
	FunctionCount int
	VariableCount int
}

// sound categorization mappings
var soundCategories = map[string][]string{
	"drums":      {"bd", "hh", "sd", "cp", "oh", "ch", "cy", "rim", "clap"},
	"percussion": {"bd", "hh", "sd", "cp", "oh", "ch", "cy", "rim", "clap"},
	"synth":      {"sawtooth", "sine", "square", "triangle", "piano", "synth"},
	"bass":       {"bass", "subbass"},
	"piano":      {"piano"},
}

// effect categorization mappings
var effectCategories = map[string]string{
	"delay":  "delay",
	"room":   "reverb",
	"lpf":    "filter",
	"hpf":    "filter",
	"crush":  "distortion",
	"gain":   "dynamics",
	"pan":    "spatial",
	"vowel":  "vocal-effects",
	"phaser": "modulation",
	"chorus": "modulation",
}

// performs full semantic analysis on Strudel code
func AnalyzeCode(code string) CodeAnalysis {
	parsed := Parse(code)

	analysis := CodeAnalysis{
		SoundTags:      analyzeSounds(code, parsed),
		EffectTags:     analyzeEffects(code, parsed),
		MusicalTags:    analyzeMusicalElements(code, parsed),
		ComplexityTags: analyzeComplexityTags(code, parsed),
		Complexity:     calculateComplexity(code, parsed),
		LineCount:      strings.Count(code, "\n") + 1,
		FunctionCount:  len(parsed.Functions),
		VariableCount:  len(parsed.Variables),
	}

	return analysis
}

// analyzeSounds categorizes sounds into semantic tags
func analyzeSounds(code string, parsed ParsedCode) []string {
	tags := make(map[string]bool)

	// check each sound against categories
	for _, sound := range parsed.Sounds {
		for category, sounds := range soundCategories {
			if contains(sounds, sound) {
				tags[category] = true
			}
		}
	}

	// check for bass pattern in code (case-insensitive)
	if regexp.MustCompile(`(?i)bass`).MatchString(code) {
		tags["bass"] = true
	}

	// convert map to slice
	return mapKeysToSlice(tags)
}

// analyzeEffects identifies audio effects used in the code
func analyzeEffects(code string, parsed ParsedCode) []string {
	tags := make(map[string]bool)

	// check each function against effect categories
	for _, function := range parsed.Functions {
		if effectTag, exists := effectCategories[function]; exists {
			tags[effectTag] = true
		}
	}

	return mapKeysToSlice(tags)
}

// analyzeMusicalElements identifies musical constructs
func analyzeMusicalElements(code string, parsed ParsedCode) []string {
	tags := make(map[string]bool)

	// note patterns indicate melody
	if len(parsed.Notes) > 0 {
		tags["melody"] = true
		tags["melodic"] = true
	}

	// scale usage
	if len(parsed.Scales) > 0 || contains(parsed.Functions, "scale") {
		tags["scales"] = true
		tags["melodic"] = true
	}

	// chord patterns (multiple notes at once indicated by comma)
	if regexp.MustCompile(`note\s*\(\s*["'][^"']*,`).MatchString(code) {
		tags["chords"] = true
		tags["harmony"] = true
	}

	// rhythm patterns
	if contains(parsed.Functions, "fast") || contains(parsed.Functions, "slow") {
		tags["rhythm"] = true
	}

	// sequences (note patterns)
	if regexp.MustCompile(`note\s*\(\s*["'][a-g0-9\s]+["']`).MatchString(code) {
		tags["sequences"] = true
	}

	return mapKeysToSlice(tags)
}

// analyzeComplexityTags generates complexity-related tags
func analyzeComplexityTags(code string, parsed ParsedCode) []string {
	tags := []string{}

	stackCount := parsed.Patterns["stack"]
	varCount := len(parsed.Variables)

	// layering tags
	if stackCount > 3 {
		tags = append(tags, "complex", "layered")
	} else if stackCount > 0 {
		tags = append(tags, "layered")
	}

	// structure tags
	if parsed.Patterns["arrange"] > 0 {
		tags = append(tags, "arranged", "structured")
	}

	// advanced code tags
	if varCount > 5 {
		tags = append(tags, "advanced")
	}

	// interactive tags
	if parsed.Patterns["slider"] > 0 {
		tags = append(tags, "interactive")
	}

	// simple/beginner-friendly tags
	if len(code) < 200 && stackCount == 0 && varCount == 0 {
		tags = append(tags, "simple", "beginner-friendly")
	}

	return tags
}

// returns a 0-10 complexity score
func calculateComplexity(code string, parsed ParsedCode) int {
	score := 0

	// base complexity from code length
	if len(code) > 500 {
		score += 3
	} else if len(code) > 200 {
		score += 2
	} else {
		score += 1
	}

	// layering complexity
	stackCount := parsed.Patterns["stack"]
	if stackCount > 3 {
		score += 3
	} else if stackCount > 0 {
		score += 2
	}

	// variable usage
	varCount := len(parsed.Variables)
	if varCount > 5 {
		score += 2
	} else if varCount > 0 {
		score += 1
	}

	// advanced patterns
	if parsed.Patterns["arrange"] > 0 {
		score += 2
	}

	// interactive elements
	if parsed.Patterns["slider"] > 0 {
		score += 1
	}

	// cap at 10
	if score > 10 {
		score = 10
	}

	return score
}

// combines analysis with existing metadata to create final tag list
func GenerateTags(analysis CodeAnalysis, category string, existingTags []string) []string {
	tags := make(map[string]bool)

	// start with existing tags
	for _, tag := range existingTags {
		if tag != "" {
			tags[strings.ToLower(tag)] = true
		}
	}

	// add category as tag
	if category != "" {
		tags[strings.ToLower(category)] = true
	}

	// add all analysis tags
	for _, tag := range analysis.SoundTags {
		tags[tag] = true
	}

	for _, tag := range analysis.EffectTags {
		tags[tag] = true
	}

	for _, tag := range analysis.MusicalTags {
		tags[tag] = true
	}

	for _, tag := range analysis.ComplexityTags {
		tags[tag] = true
	}

	return mapKeysToSlice(tags)
}

// helper functions
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func mapKeysToSlice(m map[string]bool) []string {
	result := make([]string, 0, len(m))

	for key := range m {
		result = append(result, key)
	}

	return result
}
