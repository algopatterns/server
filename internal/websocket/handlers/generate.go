package handlers

import (
	"context"
	"log"

	"github.com/algorave/server/algorave/sessions"
	"github.com/algorave/server/internal/agent"
	ws "github.com/algorave/server/internal/websocket"
)

// handles code generation request messages
func GenerateHandler(agentClient *agent.Agent, sessionRepo sessions.Repository) ws.MessageHandler {
	return func(hub *ws.Hub, client *ws.Client, msg *ws.Message) error {
		// parse payload
		var payload ws.AgentRequestPayload
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
		responsePayload := ws.AgentResponsePayload{
			Code:                response.Code,
			DocsRetrieved:       response.DocsRetrieved,
			ExamplesRetrieved:   response.ExamplesRetrieved,
			Model:               response.Model,
			IsActionable:        response.IsActionable,
			ClarifyingQuestions: response.ClarifyingQuestions,
		}

		// create response message
		responseMsg, err := ws.NewMessage(ws.TypeAgentResponse, client.SessionID, client.UserID, responsePayload)
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
