package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

// Store defines the interface for session storage backends.
// Implementations: RedisStore (primary), PostgresStore (fallback).
type Store interface {
	// Set creates or updates a session with the given TTL seconds (0 = use default idle timeout).
	Set(ctx context.Context, sessionID string, userID string, ttl int) error
	// Get retrieves the userID for a session token, extending idle TTL (sliding window).
	// Returns ("", nil) if session does not exist or is expired.
	Get(ctx context.Context, sessionID string) (string, error)
	// Delete invalidates a session (used on logout).
	Delete(ctx context.Context, sessionID string) error
	// DeleteExpired removes all expired sessions. No-op for RedisStore (TTL managed natively).
	DeleteExpired(ctx context.Context) error
}

// GenerateToken creates a cryptographically secure random session token (128 hex chars).
func GenerateToken() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
