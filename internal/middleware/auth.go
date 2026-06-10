package middleware

import (
	"context"
	"net/http"

	"go-cms/internal/session"
)

type contextKey string

// UserIDKey is the context key for the authenticated user's ID.
const UserIDKey contextKey = "userID"

// Auth protects routes by validating the session cookie.
// Valid session → injects userID into context and calls next.
// Invalid or missing session → redirects to /login.
func Auth(store session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("go_cms_session")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			userID, err := store.Get(r.Context(), cookie.Value)
			if err != nil || userID == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Inject userID into request context so handlers can read it.
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the authenticated user ID from the request context.
// Returns "" if called outside an authenticated route.
func GetUserID(r *http.Request) string {
	v, _ := r.Context().Value(UserIDKey).(string)
	return v
}
