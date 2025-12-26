package websocket

import "errors"

var (
	// errSessionNotFound indicates the session ID is invalid or doesn't exist
	ErrSessionNotFound = errors.New("session not found")

	// errUnauthorized indicates the client is not authorized for the requested action
	ErrUnauthorized = errors.New("unauthorized")

	// errInvalidMessage indicates the message format is invalid
	ErrInvalidMessage = errors.New("invalid message format")

	// errClientNotFound indicates the client ID doesn't exist in the hub
	ErrClientNotFound = errors.New("client not found")

	// errClientAlreadyRegistered indicates the client is already registered in the hub
	ErrClientAlreadyRegistered = errors.New("client already registered")

	// errSessionFull indicates the session has reached its participant limit
	ErrSessionFull = errors.New("session is full")

	// errReadOnly indicates the client doesn't have write permissions
	ErrReadOnly = errors.New("read-only access")

	// errConnectionClosed indicates the webSocket connection has been closed
	ErrConnectionClosed = errors.New("connection closed")
)
