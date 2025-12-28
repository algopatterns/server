package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/algorave/server/internal/agent"
)

func NewManager(ttl time.Duration) *Manager {
	m := &Manager{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}

	go m.cleanupExpiredSessions()

	return m
}

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (m *Manager) CreateSession() (*Session, error) {
	id, err := GenerateSessionID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &Session{
		ID:                  id,
		ConversationHistory: []agent.Message{},
		EditorState:         "",
		LastActivity:        now,
		ExpiresAt:           now.Add(m.ttl),
	}

	m.mu.Lock()
	m.sessions[id] = session
	m.mu.Unlock()

	return session, nil
}

func (m *Manager) GetSession(sessionID string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, false
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, false
	}

	return session, true
}

func (m *Manager) UpdateSession(sessionID string, history []agent.Message, editorState string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, sessionID)
		return ErrSessionExpired
	}

	now := time.Now()
	session.ConversationHistory = history
	session.EditorState = editorState
	session.LastActivity = now
	session.ExpiresAt = now.Add(m.ttl)

	return nil
}

func (m *Manager) DeleteSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, sessionID)
}

func (m *Manager) cleanupExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()

		for id, session := range m.sessions {
			if now.After(session.ExpiresAt) {
				delete(m.sessions, id)
			}
		}

		m.mu.Unlock()
	}
}

func (m *Manager) GetSessionCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.sessions)
}
