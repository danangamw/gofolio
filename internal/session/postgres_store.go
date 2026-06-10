package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"go-cms/internal/model"

	"gorm.io/gorm"
)

const (
	pgIdleTimeout     = 24 * time.Hour
	pgAbsoluteTimeout = 7 * 24 * time.Hour
)

// PostgresStore implements Store using the sessions table as the fallback backend.
// Used when REDIS_URL is not configured.
type PostgresStore struct {
	db *gorm.DB
}

// NewPostgresStore creates a PostgresStore backed by the given GORM instance.
func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// Set creates or updates a session row.
// ttl=0 uses the default idle timeout (24h).
func (s *PostgresStore) Set(ctx context.Context, sessionID, userID string, ttl int) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("session: postgres set: invalid userID %q: %w", userID, err)
	}

	idleDur := pgIdleTimeout
	if ttl > 0 {
		idleDur = time.Duration(ttl) * time.Second
	}

	now := time.Now()
	sess := model.Session{
		ID:           sessionID,
		UserID:       uid,
		ExpiresAt:    now.Add(pgAbsoluteTimeout),
		LastActiveAt: now,
	}

	// Upsert: insert or update LastActiveAt + ExpiresAt if session already exists.
	result := s.db.WithContext(ctx).
		Where(model.Session{ID: sessionID}).
		Assign(model.Session{
			LastActiveAt: now,
			ExpiresAt:    now.Add(idleDur),
		}).
		FirstOrCreate(&sess)
	return result.Error
}

// Get retrieves the userID and extends the idle TTL (sliding window).
// Returns ("", nil) if the session does not exist or is expired.
func (s *PostgresStore) Get(ctx context.Context, sessionID string) (string, error) {
	var sess model.Session
	err := s.db.WithContext(ctx).
		Where("id = ? AND expires_at > ?", sessionID, time.Now()).
		First(&sess).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("session: postgres get: %w", err)
	}

	// Extend idle TTL on each access.
	now := time.Now()
	s.db.WithContext(ctx).Model(&sess).Updates(map[string]any{
		"last_active_at": now,
		"expires_at":     now.Add(pgIdleTimeout),
	})

	return sess.UserID.String(), nil
}

// Delete removes a session row (used on logout).
func (s *PostgresStore) Delete(ctx context.Context, sessionID string) error {
	return s.db.WithContext(ctx).Delete(&model.Session{}, "id = ?", sessionID).Error
}

// DeleteExpired removes all sessions that have passed their expiry.
// Call this periodically (e.g. daily cron) to keep the sessions table clean.
func (s *PostgresStore) DeleteExpired(ctx context.Context) error {
	return s.db.WithContext(ctx).Delete(&model.Session{}, "expires_at < ?", time.Now()).Error
}
