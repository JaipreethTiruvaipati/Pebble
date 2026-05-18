// Package auth provides authentication primitives for Pebble: JWT access tokens,
// HTTP middleware, OTP login, CORS, and distributed rate limiting. This file
// defines chi/HTTP middleware that guards routes and applies in-process rate limits.
package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/jaipreeth/pebble/backend/internal/httputil"
	"golang.org/x/time/rate"
)

// contextKey is the type for request-scoped context keys to avoid collisions
// with other packages' context values.
type contextKey string

// UserIDKey is the context key under which RequireAuth stores the authenticated
// user's UUID. Handlers read r.Context().Value(UserIDKey) after the middleware chain.
const UserIDKey contextKey = "user_id"

// RequireAuth returns middleware that enforces a valid Bearer JWT on each request.
//
// Inputs: jwtManager from NewJWTManager, wired on protected route groups in api-gateway.
// Outputs: http.Handler middleware; on failure responds 401 via httputil; on success
// calls next with UserIDKey set in context from verified token claims.
//
// Expects Authorization: Bearer <access_token>. Used on portfolio, bills, investment,
// and other user-scoped API routes. Unauthenticated routes (login, health) omit this.
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

// RateLimit returns in-process token-bucket middleware keyed to the whole process.
//
// Inputs: requestsPerSecond (sustained rate) and burst (short spike allowance).
// Outputs: middleware that returns 429 when Allow() is false, otherwise passes through.
//
// Suitable for single-instance dev or coarse global caps. Production per-user limits
// use RedisRateLimit in ratelimit_redis.go so limits stay consistent across replicas.
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
