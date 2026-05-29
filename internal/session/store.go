package session

import "context"

// Store defines the interface for session storage backends.
// Implementations: RedisStore, PostgresStore.
type Store interface {
	// Set creates or updates a session.
	Set(ctx context.Context, sessionID string, userID string, ttl int) error
	// Get retrieves a userID from a session token, extending idle TTL.
	Get(ctx context.Context, sessionID string) (string, error)
	// Delete invalidates a session (logout).
	Delete(ctx context.Context, sessionID string) error
	// DeleteExpired cleans up expired sessions (only for PostgresStore).
	DeleteExpired(ctx context.Context) error
}
