package sessions

import (
	"net/http"

	"github.com/algorave/server/algorave/strudels"
	"github.com/algorave/server/internal/auth"
	"github.com/algorave/server/internal/errors"
	"github.com/algorave/server/internal/sessions"
	"github.com/gin-gonic/gin"
)

// transfers an anonymous session to an authenticated user's account
func TransferSessionHandler(sessionMgr *sessions.Manager, strudelRepo *strudels.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		var req TransferSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		session, exists := sessionMgr.GetSession(req.SessionID)
		if !exists {
			errors.NotFound(c, "session")
			return
		}

		// create a new strudel with the session's conversation history
		strudelReq := strudels.CreateStrudelRequest{
			Title:               req.Title,
			Description:         "Transferred from anonymous session",
			Code:                session.EditorState,
			ConversationHistory: session.ConversationHistory,
			IsPublic:            false,
		}

		strudel, err := strudelRepo.Create(c.Request.Context(), userID, strudelReq)
		if err != nil {
			errors.InternalError(c, "failed to create strudel", err)
			return
		}

		// delete the session after successful transfer
		sessionMgr.DeleteSession(req.SessionID)

		c.JSON(http.StatusCreated, gin.H{
			"message":    "session transferred successfully",
			"strudel":    strudel,
			"strudel_id": strudel.ID,
		})
	}
}
