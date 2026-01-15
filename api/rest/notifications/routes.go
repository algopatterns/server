package notifications

import (
	"github.com/gin-gonic/gin"

	"codeberg.org/algorave/server/internal/auth"
	"codeberg.org/algorave/server/internal/notifications"
)

func RegisterRoutes(router *gin.RouterGroup, svc *notifications.Service) {
	group := router.Group("/notifications")
	group.Use(auth.AuthMiddleware())
	{
		group.GET("", ListHandler(svc))
		group.GET("/unread-count", UnreadCountHandler(svc))
		group.POST("/:id/read", MarkReadHandler(svc))
		group.POST("/read-all", MarkAllReadHandler(svc))
	}
}
