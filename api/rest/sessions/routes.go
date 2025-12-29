package sessions

import (
	"github.com/algorave/server/algorave/sessions"
	"github.com/algorave/server/algorave/strudels"
	"github.com/algorave/server/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, sessionRepo sessions.Repository, strudelRepo *strudels.Repository) {
	router.POST("/sessions/transfer", auth.AuthMiddleware(), TransferSessionHandler(sessionRepo, strudelRepo))
}
