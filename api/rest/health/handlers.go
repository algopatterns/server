package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status:  "healthy",
		Service: "algorave",
		Version: "1.0.0",
	})
}

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
