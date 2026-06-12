package middleware

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
)

const (
	CSRFContextKey contextKey = "csrf_token"
	CSRFCookieName            = "go_cms_csrf"
)

// GetCSRFToken extracts the CSRF token from the request context.
func GetCSRFToken(ctx context.Context) string {
	if token, ok := ctx.Value(CSRFContextKey).(string); ok {
		return token
	}
	return ""
}

// CSRF validates CSRF tokens on all mutating requests (POST/PUT/DELETE/PATCH)
// using the Double Submit Cookie pattern.
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Try to read CSRF token from cookie
		cookie, err := r.Cookie(CSRFCookieName)
		if err == nil && cookie.Value != "" {
			token = cookie.Value
		} else {
			// Generate new secure random CSRF token
			b := make([]byte, 32)
			if _, err := rand.Read(b); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			token = hex.EncodeToString(b)

			// Set CSRF cookie
			http.SetCookie(w, &http.Cookie{
				Name:     CSRFCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
		}

		// Inject token into request context
		ctx := context.WithValue(r.Context(), CSRFContextKey, token)
		r = r.WithContext(ctx)

		// Validate mutating requests
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "PATCH" {
			submittedToken := r.FormValue("csrf_token")
			if submittedToken == "" {
				submittedToken = r.Header.Get("X-CSRF-Token")
			}

			if token == "" || submittedToken == "" || subtle.ConstantTimeCompare([]byte(token), []byte(submittedToken)) != 1 {
				http.Error(w, "Forbidden - CSRF token invalid or missing", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
