package sessions

import (
	"errors"
	"sync"
	"time"

	"github.com/algorave/server/internal/agent"
)

// errors
var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

// represents an anonymous user's session
type Session struct {
	ID                  string
	ConversationHistory []agent.Message
	EditorState         string
	LastActivity        time.Time
	ExpiresAt           time.Time
}

// manages anonymous user sessions in memory
type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	ttl      time.Duration
}
