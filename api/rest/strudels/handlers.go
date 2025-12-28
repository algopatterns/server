package strudels

import (
	"fmt"
	"net/http"

	"github.com/algorave/server/algorave/strudels"
	"github.com/algorave/server/internal/auth"
	"github.com/algorave/server/internal/errors"
	"github.com/gin-gonic/gin"
)

func CreateStrudelHandler(strudelRepo *strudels.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		var req strudels.CreateStrudelRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		strudel, err := strudelRepo.Create(c.Request.Context(), userID, req)
		if err != nil {
			errors.InternalError(c, "failed to create strudel", err)
			return
		}

		c.JSON(http.StatusCreated, strudel)
	}
}

func ListStrudelsHandler(strudelRepo *strudels.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		strudelsList, err := strudelRepo.List(c.Request.Context(), userID)
		if err != nil {
			errors.InternalError(c, "failed to list strudels", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"strudels": strudelsList})
	}
}

func GetStrudelHandler(strudelRepo *strudels.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		strudelID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		strudel, err := strudelRepo.Get(c.Request.Context(), strudelID, userID)
		if err != nil {
			errors.NotFound(c, "strudel")
			return
		}

		c.JSON(http.StatusOK, strudel)
	}
}

func UpdateStrudelHandler(strudelRepo *strudels.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		strudelID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		var req strudels.UpdateStrudelRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			errors.ValidationError(c, err)
			return
		}

		strudel, err := strudelRepo.Update(c.Request.Context(), strudelID, userID, req)
		if err != nil {
			errors.NotFound(c, "strudel")
			return
		}

		c.JSON(http.StatusOK, strudel)
	}
}

func DeleteStrudelHandler(strudelRepo *strudels.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := auth.GetUserID(c)
		if !exists {
			errors.Unauthorized(c, "")
			return
		}

		strudelID, ok := errors.ValidatePathUUID(c, "id")
		if !ok {
			return
		}

		err := strudelRepo.Delete(c.Request.Context(), strudelID, userID)
		if err != nil {
			errors.NotFound(c, "strudel")
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "strudel deleted"})
	}
}

func ListPublicStrudelsHandler(strudelRepo *strudels.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 50

		if l, ok := c.GetQuery("limit"); ok {
			if parsedLimit, err := parseInt(l); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}

		strudelsList, err := strudelRepo.ListPublic(c.Request.Context(), limit)
		if err != nil {
			errors.InternalError(c, "failed to list public strudels", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"strudels": strudelsList})
	}
}

func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)

	return i, err
}
