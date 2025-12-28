package websocket

import (
	"encoding/json"
	"time"
)

// message types for webSocket communication
const (
	// typeCodeUpdate is sent when a user updates the code
	TypeCodeUpdate = "code_update"

	// typeUserJoined is sent when a new user joins the session
	TypeUserJoined = "user_joined"

	// typeUserLeft is sent when a user leaves the session
	TypeUserLeft = "user_left"

	// typeAgentRequest is sent when a user requests code generation
	TypeAgentRequest = "agent_request"

	// typeAgentResponse is sent when the agent completes code generation
	TypeAgentResponse = "agent_response"

	// typeError is sent when an error occurs
	TypeError = "error"

	// typePing is sent by clients to keep the connection alive
	TypePing = "ping"

	// typePong is sent by server in response to ping
	TypePong = "pong"
)

// message represents a websocket message with typed payload
type Message struct {
	Type      string          `json:"type"`
	SessionID string          `json:"session_id"`
	ClientID  string          `json:"-"` // Internal only, not sent to clients
	UserID    string          `json:"user_id,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// contains code update information
type CodeUpdatePayload struct {
	Code        string `json:"code"`
	CursorLine  int    `json:"cursor_line,omitempty"`
	CursorCol   int    `json:"cursor_col,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// contains information about a newly joined user
type UserJoinedPayload struct {
	UserID      string `json:"user_id,omitempty"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"` // "host", "co-author", "viewer"
}

// contains information about a user who left
type UserLeftPayload struct {
	UserID      string `json:"user_id,omitempty"`
	DisplayName string `json:"display_name"`
}

// contains a code generation request
type AgentRequestPayload struct {
	UserQuery           string `json:"user_query"`
	EditorState         string `json:"editor_state,omitempty"`
	ConversationHistory []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"conversation_history,omitempty"`
}

// contains the agent's code generation response
type AgentResponsePayload struct {
	Code                string   `json:"code,omitempty"`
	DocsRetrieved       int      `json:"docs_retrieved"`
	ExamplesRetrieved   int      `json:"examples_retrieved"`
	Model               string   `json:"model"`
	IsActionable        bool     `json:"is_actionable"`
	ClarifyingQuestions []string `json:"clarifying_questions,omitempty"`
}

// contains error information (flattened to match REST API format)
type ErrorPayload struct {
	Error   string `json:"error"`             // error code (lowercase_snake_case, matches REST API)
	Message string `json:"message"`           // user-friendly message
	Details  string `json:"details,omitempty"` // optional details (sanitized in production)
}

// creates a new message with the given type and payload
func NewMessage(msgType, sessionID, userID string, payload interface{}) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:      msgType,
		SessionID: sessionID,
		UserID:    userID,
		Timestamp: time.Now(),
		Payload:   payloadBytes,
	}, nil
}

// unmarshals the payload into the provided struct
func (m *Message) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(m.Payload, v)
}
