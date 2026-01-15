package notifications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"codeberg.org/algorave/server/internal/errors"
	"codeberg.org/algorave/server/internal/notifications"
)

// @Summary List notifications
// @Description Get user's notifications with optional filtering
// @Tags notifications
// @Accept json
// @Produce json
// @Param limit query int false "Max notifications to return (1-100)" default(50)
// @Param unread query boolean false "Only return unread notifications"
// @Success 200 {object} ListResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications [get]
func ListHandler(svc *notifications.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			errors.Unauthorized(c, "authentication required")
			return
		}

		userID, ok := userIDVal.(string)
		if !ok {
			errors.Unauthorized(c, "invalid user ID")
			return
		}

		limit := 50
		if l := c.Query("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = parsed
			}
		}

		unreadOnly := c.Query("unread") == "true"
		notifs, err := svc.ListForUser(c.Request.Context(), userID, limit, unreadOnly)
		if err != nil {
			errors.InternalError(c, "failed to fetch notifications", err)
			return
		}

		unreadCount, err := svc.GetUnreadCount(c.Request.Context(), userID)
		if err != nil {
			unreadCount = 0
		}

		response := make([]NotificationResponse, 0, len(notifs))

		for _, n := range notifs {
			response = append(response, NotificationResponse{
				ID:        n.ID,
				Type:      n.Type,
				Title:     n.Title,
				Body:      n.Body,
				Data:      n.Data,
				Read:      n.Read,
				CreatedAt: n.CreatedAt,
			})
		}

		c.JSON(http.StatusOK, ListResponse{
			Notifications: response,
			UnreadCount:   unreadCount,
		})
	}
}

// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/{id}/read [post]
func MarkReadHandler(svc *notifications.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			errors.Unauthorized(c, "authentication required")
			return
		}
		userID, ok := userIDVal.(string)
		if !ok {
			errors.Unauthorized(c, "invalid user ID")
			return
		}

		notificationID := c.Param("id")
		if notificationID == "" {
			errors.BadRequest(c, "notification ID required", nil)
			return
		}

		if err := svc.MarkRead(c.Request.Context(), userID, notificationID); err != nil {
			errors.InternalError(c, "failed to mark notification as read", err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary Mark all notifications as read
// @Description Mark all user's notifications as read
// @Tags notifications
// @Accept json
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/read-all [post]
func MarkAllReadHandler(svc *notifications.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			errors.Unauthorized(c, "authentication required")
			return
		}

		userID, ok := userIDVal.(string)
		if !ok {
			errors.Unauthorized(c, "invalid user ID")
			return
		}

		if err := svc.MarkAllRead(c.Request.Context(), userID); err != nil {
			errors.InternalError(c, "failed to mark notifications as read", err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary Get unread notification count
// @Description Get the count of unread notifications for the user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} UnreadCountResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/unread-count [get]
func UnreadCountHandler(svc *notifications.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			errors.Unauthorized(c, "authentication required")
			return
		}

		userID, ok := userIDVal.(string)
		if !ok {
			errors.Unauthorized(c, "invalid user ID")
			return
		}

		count, err := svc.GetUnreadCount(c.Request.Context(), userID)
		if err != nil {
			errors.InternalError(c, "failed to get unread count", err)
			return
		}

		c.JSON(http.StatusOK, UnreadCountResponse{Count: count})
	}
}
