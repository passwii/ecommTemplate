package admin

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// minPasswordLength is the shortest password accepted for an admin account.
const minPasswordLength = 12

// ErrInvalidCredentials is returned when an email/password pair does not match.
var ErrInvalidCredentials = errors.New("admin: invalid credentials")

// HashPassword returns a bcrypt hash suitable for storage.
func HashPassword(password string) (string, error) {
	if len(password) < minPasswordLength {
		return "", fmt.Errorf("admin: password must be at least %d characters", minPasswordLength)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("admin: hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword reports whether the password matches the stored hash.
func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// NewSessionToken returns a random opaque session token and the hash to store.
// Only the hash is persisted, so a database leak does not expose live sessions.
func NewSessionToken() (token string, tokenHash string, err error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", "", fmt.Errorf("admin: generate session token: %w", err)
	}
	token = base64.RawURLEncoding.EncodeToString(buffer)
	return token, HashToken(token), nil
}

// HashToken returns the storage hash for a session token.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}

// NormalizeEmail lowercases and trims an email for consistent storage/lookup.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
