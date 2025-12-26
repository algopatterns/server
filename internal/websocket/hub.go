package websocket

import (
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to sessions
type Hub struct {
	// Registered clients by session ID and client ID
	// Structure: map[sessionID]map[clientID]*Client
	sessions map[string]map[string]*Client

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Broadcast messages to all clients in a session
	Broadcast chan *Message

	// Mutex for thread-safe access to sessions
	mu sync.RWMutex

	// Message handlers for different message types
	handlers map[string]MessageHandler

	// Flag indicating if hub is running
	running bool

	// Channel to signal shutdown
	shutdown chan struct{}
}

// MessageHandler is a function that processes a specific message type
type MessageHandler func(hub *Hub, client *Client, msg *Message) error

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		sessions:   make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 256),
		handlers:   make(map[string]MessageHandler),
		running:    false,
		shutdown:   make(chan struct{}),
	}
}

// RegisterHandler registers a handler for a specific message type
func (h *Hub) RegisterHandler(messageType string, handler MessageHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[messageType] = handler
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	h.running = true
	defer func() {
		h.running = false
	}()

	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.handleMessage(message)

		case <-h.shutdown:
			h.closeAllConnections()
			return
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Create session map if it doesn't exist
	if h.sessions[client.SessionID] == nil {
		h.sessions[client.SessionID] = make(map[string]*Client)
	}

	// Add client to session
	h.sessions[client.SessionID][client.ID] = client

	log.Printf("Client registered: %s (session: %s, role: %s, name: %s)",
		client.ID, client.SessionID, client.Role, client.DisplayName)

	// Notify other clients in the session
	userJoinedMsg, err := NewMessage(TypeUserJoined, client.SessionID, client.UserID, UserJoinedPayload{
		UserID:      client.UserID,
		DisplayName: client.DisplayName,
		Role:        client.Role,
	})
	if err == nil {
		h.broadcastToSession(client.SessionID, userJoinedMsg, client.ID)
	}
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sessionClients, exists := h.sessions[client.SessionID]
	if !exists {
		return
	}

	if _, exists := sessionClients[client.ID]; exists {
		delete(sessionClients, client.ID)
		client.Close()

		log.Printf("Client unregistered: %s (session: %s)", client.ID, client.SessionID)

		// Clean up empty sessions
		if len(sessionClients) == 0 {
			delete(h.sessions, client.SessionID)
			log.Printf("Session %s has no more clients, removed", client.SessionID)
		} else {
			// Notify other clients in the session
			userLeftMsg, err := NewMessage(TypeUserLeft, client.SessionID, client.UserID, UserLeftPayload{
				UserID:      client.UserID,
				DisplayName: client.DisplayName,
			})
			if err == nil {
				h.broadcastToSession(client.SessionID, userLeftMsg, "")
			}
		}
	}
}

// handleMessage processes an incoming message
func (h *Hub) handleMessage(msg *Message) {
	// Find the sender client by ClientID
	h.mu.RLock()
	sessionClients, exists := h.sessions[msg.SessionID]
	if !exists {
		h.mu.RUnlock()
		log.Printf("Session not found for message: %s", msg.SessionID)
		return
	}

	// Look up sender by ClientID (not UserID, to support multiple connections)
	sender, exists := sessionClients[msg.ClientID]
	h.mu.RUnlock()

	if !exists {
		log.Printf("Sender client %s not found for message in session %s", msg.ClientID, msg.SessionID)
		return
	}

	// Check if there's a registered handler for this message type
	h.mu.RLock()
	handler, exists := h.handlers[msg.Type]
	h.mu.RUnlock()

	if exists {
		// Call the handler
		if err := handler(h, sender, msg); err != nil {
			log.Printf("Handler error for message type %s: %v", msg.Type, err)
			sender.SendError("HANDLER_ERROR", "Failed to process message", err.Error())
		}
	} else {
		// Default behavior: broadcast to all clients in the session
		h.BroadcastToSession(msg.SessionID, msg, sender.ID)
	}
}

// BroadcastToSession sends a message to all clients in a session
func (h *Hub) BroadcastToSession(sessionID string, msg *Message, excludeClientID string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	h.broadcastToSession(sessionID, msg, excludeClientID)
}

// broadcastToSession is the internal broadcast function (must be called with lock held)
func (h *Hub) broadcastToSession(sessionID string, msg *Message, excludeClientID string) {
	sessionClients, exists := h.sessions[sessionID]
	if !exists {
		return
	}

	for clientID, client := range sessionClients {
		// Skip the excluded client (usually the sender)
		if clientID == excludeClientID {
			continue
		}

		if err := client.Send(msg); err != nil {
			log.Printf("Failed to send message to client %s: %v", clientID, err)
		}
	}
}

// GetSessionClients returns all clients in a session
func (h *Hub) GetSessionClients(sessionID string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sessionClients, exists := h.sessions[sessionID]
	if !exists {
		return []*Client{}
	}

	clients := make([]*Client, 0, len(sessionClients))
	for _, client := range sessionClients {
		clients = append(clients, client)
	}

	return clients
}

// GetClientCount returns the number of clients in a session
func (h *Hub) GetClientCount(sessionID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sessionClients, exists := h.sessions[sessionID]
	if !exists {
		return 0
	}

	return len(sessionClients)
}

// GetSessionCount returns the total number of active sessions
func (h *Hub) GetSessionCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.sessions)
}

// Shutdown gracefully shuts down the hub
func (h *Hub) Shutdown() {
	if h.running {
		close(h.shutdown)
	}
}

// closeAllConnections closes all client connections
func (h *Hub) closeAllConnections() {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Println("Closing all WebSocket connections...")

	for sessionID, sessionClients := range h.sessions {
		for clientID, client := range sessionClients {
			client.Close()
			log.Printf("Closed client %s in session %s", clientID, sessionID)
		}
	}

	// Clear all sessions
	h.sessions = make(map[string]map[string]*Client)
}
