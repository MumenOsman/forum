package auth

import (
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

/*
internal/auth/auth.go

Role: Security, password management, and sessions.

Responsibilities:
1. Provide functions to securely hash user passwords using bcrypt during registration.
2. Provide functions to compare a plain-text password against a hashed stored password during login.
3. Manage user sessions by generating cryptographically secure UUIDs.
4. Include middleware functions to check session cookies before allowing access to protected routes (like creating a post).
*/

// HashPassword securely hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ComparePassword compares a hashed password with a plain-text password.
func ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// GenerateSessionID securely generates a v4 UUID.
func GenerateSessionID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
