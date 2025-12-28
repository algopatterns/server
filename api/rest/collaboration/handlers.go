package collaboration

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/algorave/server/algorave/sessions"
	"github.com/algorave/server/internal/auth"
	"github.com/algorave/server/internal/errors"
	"github.com/algorave/server/internal/logger"
)

func CreateSessionHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		// parse request
		var req CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		// create session
		session, err := sessionRepo.CreateSession(c.Request.Context(), &sessions.CreateSessionRequest{
			HostUserID: userID,
			Title:      req.Title,
			Code:       req.Code,
		})
		if err != nil {
			errors.InternalError(c, "failed to create session", err)
			return
		}

		c.JSON(http.StatusCreated, CreateSessionResponse{
			ID:           session.ID,
			HostUserID:   session.HostUserID,
			Title:        session.Title,
			Code:         session.Code,
			IsActive:     session.IsActive,
			CreatedAt:    session.CreatedAt,
			LastActivity: session.LastActivity,
		})
	}
}

// retrieves a session by ID
func GetSessionHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		// get session
		session, err := sessionRepo.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			errors.SessionNotFound(c)
			return
		}

		// get participants (both authenticated and anonymous)
		participants, err := sessionRepo.ListAllParticipants(c.Request.Context(), sessionID)
		if err != nil {
			logger.ErrorErr(err, "failed to list participants",
				"session_id", sessionID,
			)
		}

		participantResponses := make([]ParticipantResponse, 0, len(participants))

		for _, p := range participants {
			participantResponses = append(participantResponses, ParticipantResponse{
				ID:          p.ID,
				UserID:      p.UserID,
				DisplayName: &p.DisplayName,
				Role:        p.Role,
				Status:      p.Status,
				JoinedAt:    p.JoinedAt,
				LeftAt:      p.LeftAt,
			})
		}

		c.JSON(http.StatusOK, SessionResponse{
			ID:           session.ID,
			HostUserID:   session.HostUserID,
			Title:        session.Title,
			Code:         session.Code,
			IsActive:     session.IsActive,
			CreatedAt:    session.CreatedAt,
			EndedAt:      session.EndedAt,
			LastActivity: session.LastActivity,
			Participants: participantResponses,
		})
	}
}

// lists all sessions for the currently authenticated user
func ListUserSessionsHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)

		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		activeOnly := c.DefaultQuery("active_only", "false") == "true"

		userSessions, err := sessionRepo.GetUserSessions(c.Request.Context(), userID, activeOnly)
		if err != nil {
			errors.InternalError(c, "failed to retrieve sessions", err)
			return
		}

		// convert to response format
		responses := make([]SessionResponse, 0, len(userSessions))

		for _, s := range userSessions {
			responses = append(responses, SessionResponse{
				ID:           s.ID,
				HostUserID:   s.HostUserID,
				Title:        s.Title,
				Code:         s.Code,
				IsActive:     s.IsActive,
				CreatedAt:    s.CreatedAt,
				EndedAt:      s.EndedAt,
				LastActivity: s.LastActivity,
			})
		}

		c.JSON(http.StatusOK, gin.H{"sessions": responses})
	}
}

// updates the strudel code in a session
func UpdateSessionCodeHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		_, err := sessionRepo.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			errors.SessionNotFound(c)
			return
		}

		// check if authenticated user is host or co-author
		participant, err := sessionRepo.GetAuthenticatedParticipant(c.Request.Context(), sessionID, userID)
		if err != nil || (participant.Role != "host" && participant.Role != "co-author") {
			errors.Forbidden(c, "you don't have permission to edit this session")
			return
		}

		// parse request
		var req UpdateSessionCodeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		// update code
		if err := sessionRepo.UpdateSessionCode(c.Request.Context(), sessionID, req.Code); err != nil {
			errors.InternalError(c, "failed to update code", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "code updated successfully", "code": req.Code})
	}
}

func EndSessionHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		session, err := sessionRepo.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			errors.SessionNotFound(c)
			return
		}

		// only authenticated user can end session
		if session.HostUserID != userID {
			errors.Forbidden(c, "only the host can end the session")
			return
		}

		if err := sessionRepo.EndSession(c.Request.Context(), sessionID); err != nil {
			errors.InternalError(c, "failed to end session", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "session ended successfully"})
	}
}

// creates an invite token for a session
func CreateInviteTokenHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		session, err := sessionRepo.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			errors.SessionNotFound(c)
			return
		}

		if session.HostUserID != userID {
			errors.Forbidden(c, "only the host can create invite tokens")
			return
		}

		var req CreateInviteTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		token, err := sessionRepo.CreateInviteToken(c.Request.Context(), &sessions.CreateInviteTokenRequest{
			SessionID: sessionID,
			Role:      req.Role,
			MaxUses:   req.MaxUses,
			ExpiresAt: req.ExpiresAt,
		})
		if err != nil {
			errors.InternalError(c, "failed to create invite token", err)
			return
		}

		c.JSON(http.StatusCreated, InviteTokenResponse{
			ID:        token.ID,
			SessionID: token.SessionID,
			Token:     token.Token,
			Role:      token.Role,
			MaxUses:   token.MaxUses,
			UsesCount: token.UsesCount,
			ExpiresAt: token.ExpiresAt,
			CreatedAt: token.CreatedAt,
		})
	}
}

func ListParticipantsHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		// get participants (both authenticated and anonymous)
		participants, err := sessionRepo.ListAllParticipants(c.Request.Context(), sessionID)
		if err != nil {
			errors.InternalError(c, "failed to retrieve participants", err)
			return
		}

		responses := make([]ParticipantResponse, 0, len(participants))

		for _, p := range participants {
			responses = append(responses, ParticipantResponse{
				ID:          p.ID,
				UserID:      p.UserID,
				DisplayName: &p.DisplayName,
				Role:        p.Role,
				Status:      p.Status,
				JoinedAt:    p.JoinedAt,
				LeftAt:      p.LeftAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{"participants": responses})
	}
}

func JoinSessionHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req JoinSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		token, err := sessionRepo.ValidateInviteToken(c.Request.Context(), req.InviteToken)
		if err != nil {
			errors.InvalidInvite(c, "")
			return
		}

		userID, _ := auth.GetUserID(c)
		displayName := req.DisplayName

		if displayName == "" {
			if userID != "" {
				displayName = "User"
			} else {
				displayName = "Anonymous"
			}
		}

		if userID != "" {
			_, err = sessionRepo.AddAuthenticatedParticipant(c.Request.Context(), token.SessionID, userID, displayName, token.Role)
			if err != nil {
				errors.InternalError(c, "failed to join session", err)
				return
			}
		} else {
			_, err = sessionRepo.AddAnonymousParticipant(c.Request.Context(), token.SessionID, displayName, token.Role)
			if err != nil {
				errors.InternalError(c, "failed to join session", err)
				return
			}
		}

		if err := sessionRepo.IncrementTokenUses(c.Request.Context(), token.ID); err != nil {
			logger.ErrorErr(err, "failed to increment token uses",
				"session_id", token.SessionID,
				"token_id", token.ID,
			)
		}

		c.JSON(http.StatusOK, JoinSessionResponse{
			SessionID:   token.SessionID,
			Role:        token.Role,
			DisplayName: displayName,
		})
	}
}

func LeaveSessionHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		participant, err := sessionRepo.GetAuthenticatedParticipant(c.Request.Context(), sessionID, userID)
		if err != nil {
			errors.NotFound(c, "you are not a participant in this session")
			return
		}

		if err := sessionRepo.MarkAuthenticatedParticipantLeft(c.Request.Context(), participant.ID); err != nil {
			errors.InternalError(c, "failed to leave session", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "successfully left session"})
	}
}

func GetSessionMessagesHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		limit := 100

		if limitStr := c.Query("limit"); limitStr != "" {
			var parsedLimit int
			if _, err := fmt.Sscanf(limitStr, "%d", &parsedLimit); err == nil {
				if parsedLimit > 0 && parsedLimit <= 1000 {
					limit = parsedLimit
				}
			}
		}

		messages, err := sessionRepo.GetMessages(c.Request.Context(), sessionID, limit)
		if err != nil {
			errors.InternalError(c, "failed to retrieve messages", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}

func RemoveParticipantHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		participantID, ok := errors.ValidatePathUUID(c, "participant_id")
		if !ok {
			return
		}

		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		session, err := sessionRepo.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			errors.SessionNotFound(c)
			return
		}

		if session.HostUserID != userID {
			errors.Forbidden(c, "only the host can remove participants")
			return
		}

		participant, err := sessionRepo.GetParticipantByID(c.Request.Context(), participantID)
		if err != nil || participant.SessionID != sessionID {
			errors.ParticipantNotFound(c)
			return
		}

		if participant.UserID != nil && *participant.UserID == userID {
			errors.InvalidOperation(c, "cannot remove yourself. use leave endpoint instead")
			return
		}

		if err := sessionRepo.RemoveParticipant(c.Request.Context(), participantID); err != nil {
			errors.InternalError(c, "failed to remove participant", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "participant removed successfully"})
	}
}

func UpdateParticipantRoleHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		participantID, ok := errors.ValidatePathUUID(c, "participant_id")
		if !ok {
			return
		}

		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		session, err := sessionRepo.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			errors.SessionNotFound(c)
			return
		}

		if session.HostUserID != userID {
			errors.Forbidden(c, "only the host can change participant roles")
			return
		}

		var req struct {
			Role string `json:"role" binding:"required,oneof=co-author viewer"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		participant, err := sessionRepo.GetParticipantByID(c.Request.Context(), participantID)
		if err != nil || participant.SessionID != sessionID {
			errors.ParticipantNotFound(c)
			return
		}

		if participant.Role == "host" {
			errors.InvalidOperation(c, "cannot change host role")
			return
		}

		if err := sessionRepo.UpdateParticipantRole(c.Request.Context(), participantID, req.Role); err != nil {
			errors.InternalError(c, "failed to update role", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "role updated successfully", "role": req.Role})
	}
}

func ListInviteTokensHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		session, err := sessionRepo.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			errors.SessionNotFound(c)
			return
		}

		if session.HostUserID != userID {
			errors.Forbidden(c, "only the host can view invite tokens")
			return
		}

		tokens, err := sessionRepo.ListInviteTokens(c.Request.Context(), sessionID)
		if err != nil {
			errors.InternalError(c, "failed to retrieve invite tokens", err)
			return
		}

		responses := make([]InviteTokenResponse, 0, len(tokens))
		for _, t := range tokens {
			responses = append(responses, InviteTokenResponse{
				ID:        t.ID,
				SessionID: t.SessionID,
				Token:     t.Token,
				Role:      t.Role,
				MaxUses:   t.MaxUses,
				UsesCount: t.UsesCount,
				ExpiresAt: t.ExpiresAt,
				CreatedAt: t.CreatedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{"tokens": responses})
	}
}

func RevokeInviteTokenHandler(sessionRepo sessions.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		tokenID, ok := errors.ValidatePathUUID(c, "token_id")
		if !ok {
			return
		}

		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		session, err := sessionRepo.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			errors.SessionNotFound(c)
			return
		}

		if session.HostUserID != userID {
			errors.Forbidden(c, "only the host can revoke invite tokens")
			return
		}

		if err := sessionRepo.RevokeInviteToken(c.Request.Context(), tokenID); err != nil {
			errors.InternalError(c, "failed to revoke invite token", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "invite token revoked successfully"})
	}
}
