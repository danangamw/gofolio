package repository

// SessionRepository handles all session database queries (fallback when Redis is unavailable).
type SessionRepository struct{}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

// TODO: Create, FindByID, Delete, DeleteExpired
