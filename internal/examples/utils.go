package examples

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// extractTags analyzes Strudel code and extracts relevant tags
// Tags help with semantic search and categorization
func extractTags(code string, category string, existingTags []string) []string {
	tags := make(map[string]bool)

	// Start with existing tags
	for _, tag := range existingTags {
		if tag != "" {
			tags[strings.ToLower(tag)] = true
		}
	}

	// Add category as tag if provided
	if category != "" {
		tags[strings.ToLower(category)] = true
	}

	// Extract sound types
	soundTags := extractSounds(code)
	for _, tag := range soundTags {
		tags[tag] = true
	}

	// Extract function/effect tags
	effectTags := extractEffects(code)
	for _, tag := range effectTags {
		tags[tag] = true
	}

	// Extract musical element tags
	musicalTags := extractMusicalElements(code)
	for _, tag := range musicalTags {
		tags[tag] = true
	}

	// Extract pattern complexity tags
	complexityTags := extractComplexity(code)
	for _, tag := range complexityTags {
		tags[tag] = true
	}

	// Convert map to slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}

	return result
}

// extractSounds finds sound types used in the code
func extractSounds(code string) []string {
	tags := make(map[string]bool)

	// Common drum sounds
	drumSounds := []string{"bd", "hh", "sd", "cp", "oh", "ch", "cy", "rim", "clap"}
	for _, sound := range drumSounds {
		if regexp.MustCompile(`sound\s*\(\s*["']` + sound + `["']`).MatchString(code) {
			tags["drums"] = true
			tags["percussion"] = true
			break
		}
	}

	// Synth/tonal sounds
	synthSounds := []string{"sawtooth", "sine", "square", "triangle", "piano", "bass", "synth"}
	for _, sound := range synthSounds {
		pattern := regexp.MustCompile(`sound\s*\(\s*["']` + sound + `["']`)
		if pattern.MatchString(code) || strings.Contains(strings.ToLower(code), sound) {
			tags["synth"] = true
			break
		}
	}

	// Check for bass
	if regexp.MustCompile(`(?i)bass`).MatchString(code) {
		tags["bass"] = true
	}

	// Check for piano
	if regexp.MustCompile(`sound\s*\(\s*["']piano["']`).MatchString(code) {
		tags["piano"] = true
	}

	// Convert to slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}
	return result
}

// extractEffects finds audio effects used in the code
func extractEffects(code string) []string {
	tags := make(map[string]bool)

	effects := map[string]string{
		`\.delay\(`:  "delay",
		`\.room\(`:   "reverb",
		`\.lpf\(`:    "filter",
		`\.hpf\(`:    "filter",
		`\.crush\(`:  "distortion",
		`\.gain\(`:   "dynamics",
		`\.pan\(`:    "spatial",
		`\.vowel\(`:  "vocal-effects",
		`\.phaser\(`: "modulation",
		`\.chorus\(`: "modulation",
	}

	for pattern, tag := range effects {
		if regexp.MustCompile(pattern).MatchString(code) {
			tags[tag] = true
		}
	}

	// Convert to slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}
	return result
}

// extractMusicalElements identifies musical constructs
func extractMusicalElements(code string) []string {
	tags := make(map[string]bool)

	// note patterns indicate melody
	if regexp.MustCompile(`note\s*\(`).MatchString(code) {
		tags["melody"] = true
		tags["melodic"] = true
	}

	// scale usage
	if regexp.MustCompile(`\.scale\s*\(`).MatchString(code) {
		tags["scales"] = true
		tags["melodic"] = true
	}

	// chord patterns (multiple notes at once)
	if regexp.MustCompile(`note\s*\(\s*["'][^"']*,`).MatchString(code) {
		tags["chords"] = true
		tags["harmony"] = true
	}

	// rhythm patterns
	if regexp.MustCompile(`\.fast\(|\.slow\(`).MatchString(code) {
		tags["rhythm"] = true
	}

	// arpeggios or sequences
	if regexp.MustCompile(`note\s*\(\s*["'][a-g0-9\s]+["']`).MatchString(code) {
		tags["sequences"] = true
	}

	// convert to slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}
	return result
}

// extractComplexity analyzes code complexity
func extractComplexity(code string) []string {
	tags := make([]string, 0)

	// count stack() calls - indicates layering
	stackCount := len(regexp.MustCompile(`stack\s*\(`).FindAllString(code, -1))
	if stackCount > 3 {
		tags = append(tags, "complex")
		tags = append(tags, "layered")
	} else if stackCount > 0 {
		tags = append(tags, "layered")
	}

	// check for arrange() - indicates structured composition
	if regexp.MustCompile(`arrange\s*\(`).MatchString(code) {
		tags = append(tags, "arranged")
		tags = append(tags, "structured")
	}

	// check for variables - indicates more advanced code
	varCount := len(regexp.MustCompile(`(?m)^(?:let|const|var)\s+`).FindAllString(code, -1))
	if varCount > 5 {
		tags = append(tags, "advanced")
	}

	// check for sliders/interactive elements
	if regexp.MustCompile(`slider\s*\(`).MatchString(code) {
		tags = append(tags, "interactive")
	}

	// simple patterns (minimal code)
	if len(code) < 200 && stackCount == 0 && varCount == 0 {
		tags = append(tags, "simple")
		tags = append(tags, "beginner-friendly")
	}

	return tags
}

// generate description creates a basic description if none exists
func generateDescription(code string, title string, category string, tags []string) string {
	// build description from available metadata
	parts := []string{}

	if category != "" {
		parts = append(parts, "A "+strings.ToLower(category)+" pattern")
	} else {
		parts = append(parts, "A Strudel pattern")
	}

	// add key features based on tags
	features := []string{}
	for _, tag := range tags {
		switch tag {
		case "drums", "percussion":
			features = append(features, "drums")
		case "melody", "melodic":
			features = append(features, "melody")
		case "bass":
			features = append(features, "bass")
		case "chords", "harmony":
			features = append(features, "chords")
		case "reverb", "delay":
			features = append(features, "effects")
		}
	}

	if len(features) > 0 {
		// deduplicate features
		uniqueFeatures := make(map[string]bool)
		for _, f := range features {
			uniqueFeatures[f] = true
		}

		featureList := make([]string, 0, len(uniqueFeatures))
		for f := range uniqueFeatures {
			featureList = append(featureList, f)
		}

		if len(featureList) == 1 {
			parts = append(parts, "featuring "+featureList[0])
		} else if len(featureList) == 2 {
			parts = append(parts, "featuring "+featureList[0]+" and "+featureList[1])
		} else if len(featureList) > 2 {
			parts = append(parts, "featuring "+strings.Join(featureList[:len(featureList)-1], ", ")+" and "+featureList[len(featureList)-1])
		}
	}

	return strings.Join(parts, " ") + "."
}

func LoadExamplesFromJSON(filePath string) ([]RawExample, error) {
	// read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read examples file: %w", err)
	}

	// parse JSON
	var rawExamples []RawExample
	if err := json.Unmarshal(data, &rawExamples); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(rawExamples) == 0 {
		return nil, fmt.Errorf("no examples found in file")
	}

	return rawExamples, nil
}

// @todo: refactor this with accurate strudel syntax
