package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/jaipreeth/pebble/backend/internal/httputil"
	"golang.org/x/time/rate"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// RequireAuth is a middleware that verifies the JWT in the Authorization header.
func RequireAuth(jwtManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				httputil.RespondError(w, http.StatusUnauthorized, "missing or invalid authorization header", "UNAUTHORIZED")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.VerifyToken(tokenStr)
			if err != nil {
				httputil.RespondError(w, http.StatusUnauthorized, "invalid or expired token", "UNAUTHORIZED")
				return
			}

			// Add the UserID to the request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RateLimit is a simple IP-based rate limiter middleware.
// In a real production environment, you'd use Redis for distributed rate limiting.
func RateLimit(requestsPerSecond float64, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				httputil.RespondError(w, http.StatusTooManyRequests, "rate limit exceeded", "RATE_LIMIT_EXCEEDED")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
