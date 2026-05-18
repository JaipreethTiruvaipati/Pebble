// Package auth provides authentication primitives for Pebble: JWT access tokens,
// HTTP middleware, OTP login, CORS, and distributed rate limiting. This file
// implements Redis-backed per-user and per-IP rate limits for horizontally scaled api-gateway.
package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/cache"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
)

// RedisRateLimit returns middleware that enforces request quotas using Redis INCR and TTL.
//
// Inputs: redis client from cache.Connect; maxRequests per window; window sliding period.
// Outputs: middleware—on Redis error fails open (allows request); on exceed returns 429
// with Retry-After; otherwise increments counter and calls next.
//
// Keys are per authenticated user (after RequireAuth) or per client IP for public routes.
// Shared across api-gateway replicas so abuse cannot bypass limits by hitting another pod.
func RedisRateLimit(redis *cache.Client, maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := rateLimitKey(r)
			ctx := r.Context()
			count, err := redis.RDB().Incr(ctx, key).Result()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if count == 1 {
				_ = redis.RDB().Expire(ctx, key, window).Err()
			}
			if int(count) > maxRequests {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
				httputil.RespondError(w, http.StatusTooManyRequests, "rate limit exceeded", "RATE_LIMIT_EXCEEDED")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// rateLimitKey builds the Redis key used for the current request's rate limit bucket.
//
// Inputs: *http.Request—uses UserIDKey from context when present (post-RequireAuth),
// else X-Forwarded-For or RemoteAddr for anonymous traffic.
// Outputs: key string "ratelimit:user:<id>" or "ratelimit:ip:<addr>".
//
// Separates authenticated users so one abusive IP cannot exhaust another user's quota.
func rateLimitKey(r *http.Request) string {
	if uid := r.Context().Value(UserIDKey); uid != nil {
		return fmt.Sprintf("ratelimit:user:%v", uid)
	}
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return fmt.Sprintf("ratelimit:ip:%s", ip)
}
