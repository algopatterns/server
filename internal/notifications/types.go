package notifications

import "time"

type Notification struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Body      *string                `json:"body,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Read      bool                   `json:"read"`
	CreatedAt time.Time              `json:"created_at"`
}

type CreateRequest struct {
	UserID string
	Type   string
	Title  string
	Body   *string
	Data   map[string]interface{}
}

// notification types
const (
	TypeAttribution = "attribution"
)
