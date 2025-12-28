package sessions

import (
	"github.com/algorave/server/algorave/anonsessions"
	"github.com/algorave/server/algorave/strudels"
	"github.com/algorave/server/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, sessionMgr *anonsessions.Manager, strudelRepo *strudels.Repository) {
	router.POST("/sessions/transfer", auth.AuthMiddleware(), TransferSessionHandler(sessionMgr, strudelRepo))
}
