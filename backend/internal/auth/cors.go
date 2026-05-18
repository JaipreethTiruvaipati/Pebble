// Package auth provides authentication primitives for Pebble: JWT access tokens,
// HTTP middleware, OTP login, CORS, and distributed rate limiting. This file
// configures Cross-Origin Resource Sharing for browser clients (web app, admin).
package auth

import (
	"net/http"
	"strings"
)

// CORS returns middleware that sets CORS response headers for allowed browser origins.
//
// Inputs: allowedOrigins—a comma-separated list from config.CORSAllowedOrigins
// (e.g. "http://localhost:5173,https://app.pebble.in"); empty or "*" allows all in dev.
// Outputs: middleware that reflects Access-Control-* headers when Origin matches,
// answers OPTIONS with 204, and otherwise delegates to next.
//
// Mounted early on the api-gateway router so SPA clients can send Authorization and
// Content-Type on cross-origin API calls with credentials.
func CORS(allowedOrigins string) func(http.Handler) http.Handler {
	origins := parseOrigins(allowedOrigins)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && originAllowed(origins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// parseOrigins splits a comma-separated origin list into trimmed entries.
//
// Inputs: raw from CORS_ALLOWED_ORIGINS; empty string or literal "*" yields nil slice.
// Outputs: slice of origin URLs, or nil meaning permissive dev default in originAllowed.
//
// Internal helper used by CORS at middleware construction time—not per request.
func parseOrigins(raw string) []string {
	if raw == "" || raw == "*" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if o := strings.TrimSpace(p); o != "" {
			out = append(out, o)
		}
	}
	return out
}

// originAllowed reports whether the request Origin header may receive CORS headers.
//
// Inputs: allowed from parseOrigins; origin from the incoming request.
// Outputs: true if allowed is empty (dev default), origin matches a listed host, or list contains "*".
//
// Prevents reflecting Access-Control-Allow-Origin for untrusted sites in production
// when CORS_ALLOWED_ORIGINS is set explicitly.
func originAllowed(allowed []string, origin string) bool {
	if len(allowed) == 0 {
		return true // dev default when CORS_ALLOWED_ORIGINS unset
	}
	for _, o := range allowed {
		if o == origin || o == "*" {
			return true
		}
	}
	return false
}
