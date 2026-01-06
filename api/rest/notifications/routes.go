package notifications

import (
	"github.com/gin-gonic/gin"

	"github.com/algrv/server/internal/auth"
	"github.com/algrv/server/internal/notifications"
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
