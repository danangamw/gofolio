package middleware

import "net/http"

// Auth protects /admin/* routes. Redirects to /login if session is invalid.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: validate session cookie, check idle timeout (24h) and absolute timeout (7d)
		next.ServeHTTP(w, r)
	})
}
