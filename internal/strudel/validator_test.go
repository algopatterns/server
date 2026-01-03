package strudel

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func getScriptDir(t *testing.T) string {
	// find the scripts directory relative to the project root
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// ind the server root (contains go.mod)
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "scripts", "validate-strudel")
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			t.Fatalf("could not find project root from %s", wd)
		}

		dir = parent
	}
}

func TestValidator_ValidCode(t *testing.T) {
	scriptDir := getScriptDir(t)

	// check if node_modules exists
	if _, err := os.Stat(filepath.Join(scriptDir, "node_modules")); os.IsNotExist(err) {
		t.Skip("node_modules not installed, run 'npm install' in scripts/validate-strudel")
	}

	v, err := NewValidator(scriptDir)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}
	defer v.Close() //nolint:errcheck // cleanup in test

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []struct {
		name string
		code string
	}{
		{"simple sound", `sound("bd sd")`},
		{"sound with effect", `sound("bd sd").fast(2)`},
		{"note pattern", `note("c3 e3 g3").sound("piano")`},
		{"stack pattern", `stack(sound("bd*4"), sound("hh*8"))`},
		{"with dollar prefix", `$: sound("bd*4")`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := v.Validate(ctx, tt.code)
			if err != nil {
				t.Fatalf("validation error: %v", err)
			}

			if !result.Valid {
				t.Errorf("expected valid code, got error: %s", result.Error)
			}
		})
	}
}

func TestValidator_InvalidCode(t *testing.T) {
	scriptDir := getScriptDir(t)

	// check if node_modules exists
	if _, err := os.Stat(filepath.Join(scriptDir, "node_modules")); os.IsNotExist(err) {
		t.Skip("node_modules not installed, run 'npm install' in scripts/validate-strudel")
	}

	v, err := NewValidator(scriptDir)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}
	defer v.Close() //nolint:errcheck // cleanup in test

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []struct {
		name string
		code string
	}{
		{"missing closing paren", `sound("bd sd"`},
		{"missing closing quote", `sound("bd sd)`},
		{"unclosed bracket", `sound("[bd sd")`},
		{"invalid syntax", `sound("bd".fast(2)`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := v.Validate(ctx, tt.code)
			if err != nil {
				t.Fatalf("validation error: %v", err)
			}
			if result.Valid {
				t.Errorf("expected invalid code, but got valid")
			}
			if result.Error == "" {
				t.Errorf("expected error message, got empty")
			}
		})
	}
}
