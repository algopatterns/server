package websocket

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
)

// returns the list of allowed origins for WebSocket connections
func GetAllowedWebSocketOrigins() []string {
	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		origins := strings.Split(envOrigins, ",")

		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}

		return origins
	}

	// default production whitelist (fallback)
	return []string{
		"https://yourdomain.com",
		"https://app.yourdomain.com",
		// add more production origins here as needed
	}
}

// checkOrigin validates the request origin based on environment
func CheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// Some clients may not send Origin header (e.g., same-origin requests)
		// In this case, we allow the connection
		return true
	}

	env := os.Getenv("ENVIRONMENT")
	if env != "production" {
		// development: allow all origins
		return true
	}

	// production: validate against whitelist
	allowedOrigins := GetAllowedWebSocketOrigins()
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	log.Printf("WebSocket origin rejected: %s (not in whitelist)", origin)
	return false
}

// generates a random client ID
func GenerateClientID() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
