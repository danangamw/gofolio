package session

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	idleTimeout = 24 * time.Hour
	keyPrefix   = "session:"
)

// RedisStore implements Store using Redis as the primary session backend.
// Session TTL is managed natively by Redis key expiration (sliding window via EXPIRE).
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore parses the Redis URL and returns a connected RedisStore.
func NewRedisStore(redisURL string) (*RedisStore, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("session: parse redis URL: %w", err)
	}
	return &RedisStore{client: redis.NewClient(opt)}, nil
}

// Set stores sessionID → userID with an expiration.
// ttl=0 uses the default idle timeout (24h).
func (s *RedisStore) Set(ctx context.Context, sessionID, userID string, ttl int) error {
	dur := idleTimeout
	if ttl > 0 {
		dur = time.Duration(ttl) * time.Second
	}
	return s.client.SetEx(ctx, keyPrefix+sessionID, userID, dur).Err()
}

// Get retrieves the userID for a session token, extending its idle TTL on each access.
// Returns ("", nil) if the session does not exist.
func (s *RedisStore) Get(ctx context.Context, sessionID string) (string, error) {
	key := keyPrefix + sessionID
	userID, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // session not found or expired — not an error
	}
	if err != nil {
		return "", fmt.Errorf("session: redis get: %w", err)
	}
	// Sliding window: extend idle TTL on each active request
	s.client.Expire(ctx, key, idleTimeout)
	return userID, nil
}

// Delete removes a session from Redis (logout).
func (s *RedisStore) Delete(ctx context.Context, sessionID string) error {
	return s.client.Del(ctx, keyPrefix+sessionID).Err()
}

// DeleteExpired is a no-op for Redis — TTL expiration is handled natively.
func (s *RedisStore) DeleteExpired(_ context.Context) error {
	return nil
}
