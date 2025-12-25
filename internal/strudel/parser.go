package strudel

import (
	"regexp"
	"strings"
)

// compiled regex patterns (shared across all functions)
var (
	// sound extraction: sound("bd") or s("bd")
	soundPattern = regexp.MustCompile(`(?:sound|s)\s*\(\s*["']([^"']+)["']`)

	// note extraction: note("c e g")
	notePattern = regexp.MustCompile(`note\s*\(\s*["']([^"']+)["']`)

	// function calls: .fast(2), .slow(4), .stack()
	functionPattern = regexp.MustCompile(`\.(\w+)\s*\(`)

	// variable declarations: let x = ..., const y = ...
	variablePattern = regexp.MustCompile(`(?:let|const|var)\s+(\w+)\s*=`)

	// scale/mode: scale("minor"), mode("dorian")
	scalePattern = regexp.MustCompile(`(?:scale|mode)\s*\(\s*["'](\w+)["']`)
)

// contains all extracted elements from Strudel code
type ParsedCode struct {
	Sounds    []string       // sound sample names: ["bd", "hh", "sd"]
	Notes     []string       // note names: ["c", "e", "g"]
	Functions []string       // function names: ["fast", "slow", "stack"]
	Variables []string       // variable names: ["pat1", "rhythm"]
	Scales    []string       // scale/mode names: ["minor", "dorian"]
	Patterns  map[string]int // pattern counts: {"stack": 2, "arrange": 1}
}

// extracts all elements from Strudel code
func Parse(code string) ParsedCode {
	return ParsedCode{
		Sounds:    ExtractSounds(code),
		Notes:     ExtractNotes(code),
		Functions: ExtractFunctions(code),
		Variables: ExtractVariables(code),
		Scales:    ExtractScales(code),
		Patterns:  CountPatterns(code),
	}
}

// extracts sound sample names from sound() calls
// example: sound("bd hh sd") → ["bd", "hh", "sd"]
func ExtractSounds(code string) []string {
	sounds := []string{}

	matches := soundPattern.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			// split space-separated sounds: "bd hh sd" → ["bd", "hh", "sd"]
			soundList := strings.Fields(match[1])
			sounds = append(sounds, soundList...)
		}
	}

	return sounds
}

// extractNotes extracts note names from note() calls
// example: note("c e g") → ["c", "e", "g"]
func ExtractNotes(code string) []string {
	notes := []string{}

	matches := notePattern.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			// Split space-separated notes
			noteList := strings.Fields(match[1])
			notes = append(notes, noteList...)
		}
	}

	return notes
}

// extractFunctions extracts function/method names from .func() calls
// example: .fast(2).slow(4) → ["fast", "slow"]
func ExtractFunctions(code string) []string {
	functions := []string{}

	matches := functionPattern.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			functions = append(functions, match[1])
		}
	}

	return functions
}

// extractVariables extracts variable names from declarations
// example: let pat1 = sound("bd") → ["pat1"]
func ExtractVariables(code string) []string {
	variables := []string{}

	matches := variablePattern.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables = append(variables, match[1])
		}
	}

	return variables
}

// extractScales extracts scale/mode names
// example: scale("minor") → ["minor"]
func ExtractScales(code string) []string {
	scales := []string{}

	matches := scalePattern.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			scales = append(scales, match[1])
		}
	}

	return scales
}

// counts occurrences of specific patterns
func CountPatterns(code string) map[string]int {
	patterns := make(map[string]int)

	// common patterns to count
	patternRegexes := map[string]*regexp.Regexp{
		"stack":   regexp.MustCompile(`\.stack\s*\(`),
		"every":   regexp.MustCompile(`\.every\s*\(`),
		"arrange": regexp.MustCompile(`(?:^|\W)arrange\s*\(`), // match arrange() as function or method
		"slider":  regexp.MustCompile(`slider\s*\(`),
	}

	for name, regex := range patternRegexes {
		matches := regex.FindAllString(code, -1)
		patterns[name] = len(matches)
	}

	return patterns
}

// counts occurrences of a specific pattern
func CountPattern(code string, pattern string) int {
	regex := regexp.MustCompile(`\.` + pattern + `\s*\(`)
	matches := regex.FindAllString(code, -1)

	return len(matches)
}
