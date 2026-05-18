// Package auth provides authentication primitives for Pebble: JWT access tokens,
// HTTP middleware, OTP login, CORS, and distributed rate limiting. This file
// will host refresh-token rotation (long-lived cookie or body token, Redis denylist,
// new access JWT) using config.JWTRefreshExpiryDays—not yet implemented.
package auth

// TODO: implement refresh
