package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// Auth is a simple JWT middleware stub.
// Replace with real JWT validation in production.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			// Allow unauthenticated in dev — set a default user
			ctx := context.WithValue(r.Context(), UserIDKey, "dev-user")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Strip "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")

		// TODO: validate JWT, extract user ID
		userID := "user-from-token"
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) string {
	v, _ := ctx.Value(UserIDKey).(string)
	return v
}
