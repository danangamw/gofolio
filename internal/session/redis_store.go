package session

// RedisStore implements Store using Redis as the session backend.
// Session TTL is managed natively by Redis key expiration.
type RedisStore struct{}

// TODO: implement Set, Get, Delete using go-redis client
