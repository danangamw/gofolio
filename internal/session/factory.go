package session

import (
	"fmt"

	"go-cms/internal/database"
)

// NewStore selects the session backend automatically:
//   - If redisURL is non-empty → RedisStore (primary, recommended)
//   - Otherwise → PostgresStore (fallback, no extra infra needed)
func NewStore(redisURL string, db database.Service) (Store, error) {
	if redisURL != "" {
		store, err := NewRedisStore(redisURL)
		if err != nil {
			return nil, fmt.Errorf("session: init redis store: %w", err)
		}
		return store, nil
	}
	return NewPostgresStore(db.GetDB()), nil
}
