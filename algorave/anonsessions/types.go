package anonsessions

import (
	"sync"
	"time"

	"github.com/algorave/server/algorave/strudels"
	"github.com/algorave/server/internal/agent"
)

// represents an in-memory session for anonymous users
type AnonymousSession struct {
	ID                  string                          `json:"id"`
	EditorState         string                          `json:"editor_state"`
	ConversationHistory strudels.ConversationHistory    `json:"conversation_history"`
	CreatedAt           time.Time                       `json:"created_at"`
	LastActivity        time.Time                       `json:"last_activity"`
	mu                  sync.RWMutex                    `json:"-"`
}

// updates the editor state and last activity time
func (s *AnonymousSession) UpdateEditorState(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.EditorState = code
	s.LastActivity = time.Now()
}

// safely retrieves the editor state
func (s *AnonymousSession) GetEditorState() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.EditorState
}

// adds a message to the conversation history
func (s *AnonymousSession) AddMessage(role, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ConversationHistory = append(s.ConversationHistory, agent.Message{
		Role:    role,
		Content: content,
	})
	s.LastActivity = time.Now()
}

// safely retrieves the conversation history
func (s *AnonymousSession) GetConversationHistory() strudels.ConversationHistory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ConversationHistory
}

// updates the last activity time
func (s *AnonymousSession) Touch() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastActivity = time.Now()
}
