package middleware

import "net/http"

// RateLimiter limits requests per IP.
// /login: max 10 req/min. Public routes: max 60 req/min.
func RateLimiter(limit int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement token bucket or sliding window per IP
			next.ServeHTTP(w, r)
		})
	}
}
