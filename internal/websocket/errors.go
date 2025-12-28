package websocket

import (
	"os"
	"strings"
)

// sanitizes error details for production
func sanitizeErrorString(errMsg string) string {
	if errMsg == "" {
		return ""
	}

	env := os.Getenv("ENVIRONMENT")
	if env != "production" {
		return errMsg
	}

	if strings.Contains(errMsg, "database") || strings.Contains(errMsg, "sql") {
		return "database operation failed"
	}

	if strings.Contains(errMsg, "connection") || strings.Contains(errMsg, "network") {
		return "connection error occurred"
	}

	if strings.Contains(errMsg, "timeout") {
		return "request timed out"
	}

	if strings.Contains(errMsg, "permission") || strings.Contains(errMsg, "unauthorized") {
		return "permission denied"
	}

	if strings.Contains(errMsg, "not found") {
		return "resource not found"
	}

	return "an error occurred"
}
