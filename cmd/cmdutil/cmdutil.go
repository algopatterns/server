package cmdutil

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func ParseFlags() (*string, *bool) {
	docsPath := flag.String("docs", "./docs/strudel", "path to documentation directory")
	clearExisting := flag.Bool("clear", false, "clear existing chunks before ingesting")

	flag.Parse()

	return docsPath, clearExisting
}

func LoadEnvironmentVariables() (*string, *string, error) {
	// Try to load .env file, but don't fail if it doesn't exist (CI uses env vars)
	if err := godotenv.Load(); err != nil {
		// Not fatal - environment variables might be set directly
		_ = err
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	supabaseConnStr := os.Getenv("SUPABASE_CONNECTION_STRING")

	if openaiKey == "" {
		return nil, nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	if supabaseConnStr == "" {
		return nil, nil, fmt.Errorf("SUPABASE_CONNECTION_STRING environment variable is required")
	}

	return &openaiKey, &supabaseConnStr, nil
}
