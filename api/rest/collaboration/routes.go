package collaboration

import (
	"github.com/gin-gonic/gin"

	"github.com/algorave/server/algorave/sessions"
	"github.com/algorave/server/internal/auth"
)

// RegisterRoutes registers all collaboration/session routes
func RegisterRoutes(router *gin.RouterGroup, sessionRepo sessions.Repository) {
	// Session management (authenticated)
	router.POST("/sessions", auth.AuthMiddleware(), CreateSessionHandler(sessionRepo))
	router.GET("/sessions", auth.AuthMiddleware(), ListUserSessionsHandler(sessionRepo))
	router.GET("/sessions/:id", GetSessionHandler(sessionRepo))
	router.PUT("/sessions/:id", auth.AuthMiddleware(), UpdateSessionCodeHandler(sessionRepo))
	router.DELETE("/sessions/:id", auth.AuthMiddleware(), EndSessionHandler(sessionRepo))

	// Invite tokens (host only)
	router.POST("/sessions/:id/invite", auth.AuthMiddleware(), CreateInviteTokenHandler(sessionRepo))

	// Participants
	router.GET("/sessions/:id/participants", ListParticipantsHandler(sessionRepo))

	// Join session (optional auth)
	router.POST("/sessions/join", auth.OptionalAuthMiddleware(), JoinSessionHandler(sessionRepo))
}
