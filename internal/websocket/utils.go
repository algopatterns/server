package websocket

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

func GetAllowedWebSocketOrigins() []string {
	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		origins := strings.Split(envOrigins, ",")

		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}

		return origins
	}

	return []string{}
}

func CheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	if origin == "" {
		// allow no origin header in development
		env := os.Getenv("ENVIRONMENT")

		if env != "production" {
			return true
		}

		log.Printf("webSocket connection with no Origin header")
		return false
	}

	env := os.Getenv("ENVIRONMENT")
	if env != "production" {

		return true
	}

	// production: validate against allowed origins
	allowedOrigins := GetAllowedWebSocketOrigins()
	if len(allowedOrigins) == 0 {
		log.Printf("webSocket origin rejected: %s (ALLOWED_ORIGINS not configured in production)", origin)
		return false
	}

	if slices.Contains(allowedOrigins, origin) {
		return true
	}

	log.Printf("webSocket origin rejected: %s (not in allowed origins: %v)", origin, allowedOrigins)
	return false
}

func GenerateClientID() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
