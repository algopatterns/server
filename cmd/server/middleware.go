package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

// configures cors-origin resource sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigin := os.Getenv("CORS_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:3000"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
