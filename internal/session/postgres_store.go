package session

// PostgresStore implements Store using the sessions table as a fallback backend.
// Used when REDIS_URL is not configured.
type PostgresStore struct{}

// TODO: implement Set, Get, Delete, DeleteExpired using GORM
