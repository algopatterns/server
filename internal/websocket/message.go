package websocket

import (
	"encoding/json"
	"time"
)

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
