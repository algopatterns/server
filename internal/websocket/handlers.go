package websocket

import (
	"context"
	"log"

	"github.com/algorave/server/algorave/sessions"
	"github.com/algorave/server/internal/agent"
)

// CodeUpdateHandler handles code update messages
func CodeUpdateHandler(sessionRepo sessions.Repository) MessageHandler {
	return func(hub *Hub, client *Client, msg *Message) error {
		// check if client has write permissions
		if !client.CanWrite() {
			client.SendError("FORBIDDEN", "You don't have permission to edit code", "")
			return ErrReadOnly
		}

		// parse payload
		var payload CodeUpdatePayload
		if err := msg.UnmarshalPayload(&payload); err != nil {
			client.SendError("INVALID_PAYLOAD", "Failed to parse code update", err.Error())
			return err
		}

		// update session code in database
		ctx := context.Background()
		if err := sessionRepo.UpdateSessionCode(ctx, client.SessionID, payload.Code); err != nil {
			log.Printf("Failed to update session code: %v", err)
			client.SendError("DATABASE_ERROR", "Failed to save code update", err.Error())
			return err
		}

		// add display name to payload
		payload.DisplayName = client.DisplayName

		// create new message with updated payload
		broadcastMsg, err := NewMessage(TypeCodeUpdate, client.SessionID, client.UserID, payload)
		if err != nil {
			log.Printf("Failed to create broadcast message: %v", err)
			return err
		}

		// broadcast to all other clients in the session
		hub.BroadcastToSession(client.SessionID, broadcastMsg, client.ID)

		log.Printf("Code updated by %s in session %s", client.DisplayName, client.SessionID)

		return nil
	}
}

// GenerateHandler handles code generation request messages
func GenerateHandler(agentClient *agent.Agent, sessionRepo sessions.Repository) MessageHandler {
	return func(hub *Hub, client *Client, msg *Message) error {
		// parse payload
		var payload AgentRequestPayload
		if err := msg.UnmarshalPayload(&payload); err != nil {
			client.SendError("INVALID_PAYLOAD", "Failed to parse generation request", err.Error())
			return err
		}

		// convert conversation history to agent.Message format
		conversationHistory := make([]agent.Message, 0, len(payload.ConversationHistory))
		for _, msg := range payload.ConversationHistory {
			conversationHistory = append(conversationHistory, agent.Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}

		// create agent request
		agentReq := agent.GenerateRequest{
			UserQuery:           payload.UserQuery,
			EditorState:         payload.EditorState,
			ConversationHistory: conversationHistory,
		}

		// generate code using agent
		ctx := context.Background()
		response, err := agentClient.Generate(ctx, agentReq)
		if err != nil {
			log.Printf("Failed to generate code: %v", err)
			client.SendError("GENERATION_ERROR", "Failed to generate code", err.Error())
			return err
		}

		// save messages to session history
		if payload.UserQuery != "" {
			_, err := sessionRepo.AddMessage(ctx, client.SessionID, client.UserID, "user", payload.UserQuery)
			if err != nil {
				log.Printf("Failed to save user message: %v", err)
			}
		}

		if response.Code != "" {
			_, err := sessionRepo.AddMessage(ctx, client.SessionID, "", "assistant", response.Code)
			if err != nil {
				log.Printf("Failed to save assistant message: %v", err)
			}
		}

		// update session code if generation was successful
		if response.IsActionable && response.Code != "" {
			if err := sessionRepo.UpdateSessionCode(ctx, client.SessionID, response.Code); err != nil {
				log.Printf("Failed to update session code: %v", err)
			}
		}

		// create response payload
		responsePayload := AgentResponsePayload{
			Code:                response.Code,
			DocsRetrieved:       response.DocsRetrieved,
			ExamplesRetrieved:   response.ExamplesRetrieved,
			Model:               response.Model,
			IsActionable:        response.IsActionable,
			ClarifyingQuestions: response.ClarifyingQuestions,
		}

		// create response message
		responseMsg, err := NewMessage(TypeAgentResponse, client.SessionID, client.UserID, responsePayload)
		if err != nil {
			log.Printf("Failed to create response message: %v", err)
			return err
		}

		// send response to all clients in the session (including requester)
		hub.BroadcastToSession(client.SessionID, responseMsg, "")

		// update last activity
		if err := sessionRepo.UpdateLastActivity(ctx, client.SessionID); err != nil {
			log.Printf("Failed to update last activity: %v", err)
		}

		log.Printf("Code generated for session %s by %s", client.SessionID, client.DisplayName)

		return nil
	}
}
