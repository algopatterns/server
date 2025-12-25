package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIKey          string
	SupabaseConnString string
}

func LoadEnvironmentVariables() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		_ = err
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	supabaseConnStr := os.Getenv("SUPABASE_CONNECTION_STRING")

	if openaiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	if supabaseConnStr == "" {
		return nil, fmt.Errorf("SUPABASE_CONNECTION_STRING environment variable is required")
	}

	return &Config{
		OpenAIKey:          openaiKey,
		SupabaseConnString: supabaseConnStr,
	}, nil
}
