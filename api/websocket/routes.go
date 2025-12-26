package websocket

import (
	"github.com/gin-gonic/gin"

	"github.com/algorave/server/algorave/sessions"
	ws "github.com/algorave/server/internal/websocket"
)

func RegisterRoutes(router *gin.RouterGroup, hub *ws.Hub, sessionRepo sessions.Repository) {
	router.GET("/ws", WebSocketHandler(hub, sessionRepo))
}
