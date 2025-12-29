package users

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// daily generation limits by tier/type
const (
	DailyLimitAnonymous = 50   // anonymous users: 50/day
	DailyLimitFree      = 100  // free tier: 100/day
	DailyLimitPro       = 1000 // pro tier: 1000/day
	DailyLimitBYOK      = -1   // BYOK: unlimited (using own keys)
)

// per-minute generation limits by tier (for future use)
// currently all users share the same per-minute limit (10/min) defined in websocket/types.go
const (
	MinuteLimitDefault = 10 // default for all users
	MinuteLimitPro     = 20 // pro tier: higher burst capacity
	MinuteLimitBYOK    = 30 // BYOK: highest burst capacity
)

// handles user database operations
type Repository struct {
	db *pgxpool.Pool
}

// represents an authenticated user in the system
type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"-"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatar_url"`
	Tier       string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// contains data for updating a user's profile
type UpdateProfileRequest struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// contains data for logging a generation request
type UsageLogRequest struct {
	UserID       *string // nil for anonymous
	SessionID    string  // for anonymous users
	Provider     string  // "anthropic", "openai"
	Model        string  // model name
	InputTokens  int     // estimated input tokens
	OutputTokens int     // estimated output tokens
	IsBYOK       bool    // true if user provided own API key
}

// result of a rate limit check
type RateLimitResult struct {
	Allowed   bool
	Current   int
	Limit     int
	Remaining int
}
