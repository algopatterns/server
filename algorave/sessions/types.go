package sessions

import (
	"time"
)

// Message type constants for session messages
const (
	MessageTypeUserPrompt = "user"
	MessageTypeAIResponse = "assistant"
	MessageTypeChat       = "user"
)

// represents a collaborative coding session
type Session struct {
	ID           string     `json:"id"`
	HostUserID   string     `json:"host_user_id"`
	Title        string     `json:"title"`
	Code         string     `json:"code"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	LastActivity time.Time  `json:"last_activity"`
}

// represents an authenticated user in a session
type Participant struct {
	ID          string     `json:"id"`
	SessionID   string     `json:"session_id"`
	UserID      string     `json:"user_id"`
	DisplayName string     `json:"display_name"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	JoinedAt    time.Time  `json:"joined_at"`
	LeftAt      *time.Time `json:"left_at,omitempty"`
}

// represents an anonymous user in a session
type AnonymousParticipant struct {
	ID          string     `json:"id"`
	SessionID   string     `json:"session_id"`
	DisplayName string     `json:"display_name"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	JoinedAt    time.Time  `json:"joined_at"`
	LeftAt      *time.Time `json:"left_at,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at"`
}

// represents either an authenticated or anonymous participant
type CombinedParticipant struct {
	ID          string     `json:"id"`
	SessionID   string     `json:"session_id"`
	UserID      *string    `json:"user_id,omitempty"`
	DisplayName string     `json:"display_name"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	JoinedAt    time.Time  `json:"joined_at"`
	LeftAt      *time.Time `json:"left_at,omitempty"`
}

// represents a session invite token
type InviteToken struct {
	ID        string     `json:"id"`
	SessionID string     `json:"session_id"`
	Token     string     `json:"token"`
	Role      string     `json:"role"`
	MaxUses   *int       `json:"max_uses,omitempty"`
	UsesCount int        `json:"uses_count"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// represents a chat message in a session
type Message struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"sessionID"`
	UserID      *string   `json:"userID,omitempty"`
	Role        string    `json:"role"` // user, assistant
	MessageType string    `json:"messageType"` // MessageTypeUserPrompt, MessageTypeAIResponse, MessageTypeChat
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
}

// contains data for creating a session
type CreateSessionRequest struct {
	HostUserID string `json:"host_user_id"`
	Title      string `json:"title"`
	Code       string `json:"code"`
}

// contains data for creating an invite token
type CreateInviteTokenRequest struct {
	SessionID string     `json:"session_id"`
	Role      string     `json:"role"`
	MaxUses   *int       `json:"max_uses,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
