package llm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestNewLLMWithConfig_ValidateProviders(t *testing.T) {
	// ensure invalid providers are rejected
	tests := []struct {
		name          string
		config        *Config
		expectError   bool
		errorContains string
	}{
		{
			name:          "nil config",
			config:        nil,
			expectError:   true,
			errorContains: "config cannot be nil",
		},
		{
			name: "unsupported transformer provider",
			config: &Config{
				TransformerProvider: Provider("invalid"),
				TransformerAPIKey:   "key",
				GeneratorProvider:   ProviderAnthropic,
				GeneratorAPIKey:     "key",
				EmbedderProvider:    ProviderOpenAI,
				EmbedderAPIKey:      "key",
			},
			expectError:   true,
			errorContains: "unsupported transformer provider",
		},
		{
			name: "unsupported generator provider",
			config: &Config{
				TransformerProvider: ProviderAnthropic,
				TransformerAPIKey:   "key",
				GeneratorProvider:   Provider("invalid"),
				GeneratorAPIKey:     "key",
				EmbedderProvider:    ProviderOpenAI,
				EmbedderAPIKey:      "key",
			},
			expectError:   true,
			errorContains: "unsupported generator provider",
		},
		{
			name: "unsupported embedder provider",
			config: &Config{
				TransformerProvider: ProviderAnthropic,
				TransformerAPIKey:   "key",
				GeneratorProvider:   ProviderAnthropic,
				GeneratorAPIKey:     "key",
				EmbedderProvider:    Provider("invalid"),
				EmbedderAPIKey:      "key",
			},
			expectError:   true,
			errorContains: "unsupported embedder provider",
		},
		{
			name: "valid anthropic config",
			config: &Config{
				TransformerProvider: ProviderAnthropic,
				TransformerAPIKey:   "key",
				GeneratorProvider:   ProviderAnthropic,
				GeneratorAPIKey:     "key",
				EmbedderProvider:    ProviderOpenAI,
				EmbedderAPIKey:      "key",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLLMWithConfig(context.Background(), tt.config)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRateLimiter_BurstCapacity(t *testing.T) {
	testLimiter := rate.NewLimiter(1, 2)
	ctx := context.Background()

	// first two requests should succeed immediately (burst capacity)
	err1 := testLimiter.Wait(ctx)
	err2 := testLimiter.Wait(ctx)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// third request would need to wait ~1 second
	ctx3, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	err3 := testLimiter.Wait(ctx3)

	// should timeout because rate limit requires waiting
	require.Error(t, err3)
}

func TestAnthropicConfig_Defaults(t *testing.T) {
	transformer := NewAnthropicTransformer(AnthropicConfig{
		APIKey: "test-key",
		Model:  "claude-3-haiku-20240307",
	})

	assert.Equal(t, defaultMaxTokens, transformer.config.MaxTokens)
	assert.Equal(t, float32(defaultTemperature), transformer.config.Temperature)
}

func TestOpenAIConfig_Defaults(t *testing.T) {
	embedder := NewOpenAIEmbedder(OpenAIConfig{
		APIKey: "test-key",
	})

	assert.Equal(t, defaultOpenAIModel, embedder.config.Model)

	generator := NewOpenAIGenerator(OpenAIConfig{
		APIKey: "test-key",
	})

	assert.Equal(t, defaultOpenAIChatModel, generator.config.Model)
}

// test if context cancellation is respected by rate limiter
func TestContextCancellation_RateLimiter(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel

	// rate limiter should return error when context is cancelled
	limiter := rate.NewLimiter(1, 1)
	err := limiter.Wait(ctx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

// test if empty input is rejected before making API call
func TestEmbeddings_InputValidation(t *testing.T) {

	embedder := NewOpenAIEmbedder(OpenAIConfig{
		APIKey: "test-key",
		Model:  "text-embedding-3-small",
	})

	ctx := context.Background()
	_, err := embedder.GenerateEmbeddings(ctx, []string{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no texts provided")
}
