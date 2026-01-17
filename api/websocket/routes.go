package websocket

import (
	"github.com/gin-gonic/gin"

	"codeberg.org/algorave/server/algorave/sessions"
	"codeberg.org/algorave/server/algorave/users"
	ws "codeberg.org/algorave/server/internal/websocket"
)

func RegisterRoutes(router *gin.RouterGroup, hub *ws.Hub, sessionRepo sessions.Repository, userRepo *users.Repository) {
	router.GET("/ws", WebSocketHandler(hub, sessionRepo, userRepo))
}
