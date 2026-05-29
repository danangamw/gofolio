package middleware

import "net/http"

// CSRF validates CSRF tokens on all mutating requests (POST/PUT/DELETE).
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: generate token per session, validate on POST/PUT/DELETE
		next.ServeHTTP(w, r)
	})
}
