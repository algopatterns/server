package handlers

import (
	"context"
	"log"

	"github.com/algorave/server/algorave/sessions"
	ws "github.com/algorave/server/internal/websocket"
)

// handles code update messages
func CodeUpdateHandler(sessionRepo sessions.Repository) ws.MessageHandler {
	return func(hub *ws.Hub, client *ws.Client, msg *ws.Message) error {
		// check if client has write permissions
		if !client.CanWrite() {
			client.SendError("FORBIDDEN", "You don't have permission to edit code", "")
			return ws.ErrReadOnly
		}

		// parse payload
		var payload ws.CodeUpdatePayload
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
		broadcastMsg, err := ws.NewMessage(ws.TypeCodeUpdate, client.SessionID, client.UserID, payload)
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
