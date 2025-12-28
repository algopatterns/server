package auth

import (
	"github.com/algorave/server/algorave/users"
	"github.com/golang-jwt/jwt/v5"
)

// User is an alias for users.User for backward compatibility
type User = users.User

// represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
