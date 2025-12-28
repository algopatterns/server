package generate

import (
	"context"
	"net/http"

	"github.com/algorave/server/algorave/strudels"
	"github.com/algorave/server/internal/agent"
	"github.com/algorave/server/internal/errors"
	"github.com/algorave/server/internal/logger"
	"github.com/algorave/server/internal/sessions"
	"github.com/gin-gonic/gin"
)

type StrudelGetter interface {
	Get(ctx context.Context, strudelID, userID string) (*strudels.Strudel, error)
}

// creates a handler for code generation
func Handler(agentClient *agent.Agent, strudelRepo StrudelGetter, sessionMgr *sessions.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request

		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		// check if user is authenticated
		userID, isAuthenticated := c.Get("user_id")

		// use conversation history from request by default
		conversationHistory := req.ConversationHistory
		editorState := req.EditorState
		sessionID := req.SessionID

		// validate optional UUIDs if provided
		if !errors.ValidateUUID(c, req.StrudelID, "strudel") {
			return
		}

		if !errors.ValidateUUID(c, req.SessionID, "session") {
			return
		}

		// priority 1: if strudel_id is provided and user is authenticated, load history from strudel
		if req.StrudelID != "" && isAuthenticated {
			strudel, err := strudelRepo.Get(c.Request.Context(), req.StrudelID, userID.(string))
			if err != nil {
				logger.Warn("failed to load strudel",
					"strudel_id", req.StrudelID,
					"user_id", userID,
					"error", err,
				)
			} else {
				conversationHistory = strudel.ConversationHistory
				editorState = strudel.Code
			}
		} else if !isAuthenticated {
			// priority 2: for anonymous users, check for session
			if req.SessionID != "" {
				// try to load existing session
				session, exists := sessionMgr.GetSession(req.SessionID)

				if exists {
					conversationHistory = session.ConversationHistory
					editorState = session.EditorState
					sessionID = session.ID
				} else {
					// session expired or invalid, create new one
					newSession, err := sessionMgr.CreateSession()
					if err != nil {
						logger.ErrorErr(err, "failed to create session")
					} else {
						sessionID = newSession.ID
					}
				}
			} else {
				// no session_id provided, create new session
				newSession, err := sessionMgr.CreateSession()
				if err != nil {
					logger.ErrorErr(err, "failed to create session")
				} else {
					sessionID = newSession.ID
				}
			}
		}

		resp, err := agentClient.Generate(c.Request.Context(), agent.GenerateRequest{
			UserQuery:           req.UserQuery,
			EditorState:         editorState,
			ConversationHistory: conversationHistory,
		})

		if err != nil {
			errors.InternalError(c, "failed to generate code", err)
			return
		}

		// update session for anonymous users
		if !isAuthenticated && sessionID != "" {
			// append new conversation turn
			updatedHistory := append(conversationHistory,
				agent.Message{Role: "user", Content: req.UserQuery},
				agent.Message{Role: "assistant", Content: resp.Code},
			)

			err := sessionMgr.UpdateSession(sessionID, updatedHistory, resp.Code)
			if err != nil {
				logger.ErrorErr(err, "failed to update session",
					"session_id", sessionID,
				)
			}
		}

		c.JSON(http.StatusOK, Response{
			Code:                resp.Code,
			DocsRetrieved:       resp.DocsRetrieved,
			ExamplesRetrieved:   resp.ExamplesRetrieved,
			Model:               resp.Model,
			IsActionable:        resp.IsActionable,
			ClarifyingQuestions: resp.ClarifyingQuestions,
			SessionID:           sessionID, // return session ID when available (required for anonymous users, optional for authenticated)
		})
	}
}
