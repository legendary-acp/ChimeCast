package middleware

import (
	"context"
	"net/http"

	"github.com/legendary-acp/chimecast/internal/session"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(sessionManager *session.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := sessionManager.GetSession(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add both userID and username to context
			ctx := context.WithValue(r.Context(), "userID", session.UserID)
			ctx = context.WithValue(ctx, "userName", session.UserName)

			// Create new request with the updated context
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
