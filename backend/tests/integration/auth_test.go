//go:build integration

// Package integration holds end-to-end tests that require PostgreSQL, Redis, and RabbitMQ.
// Build with: go test ./backend/tests/... -tags integration -v
package integration

import (
	"testing"
)

// TestAuthFlow verifies signup → OTP → login → JWT-protected /users/me (planned; uses testcontainers).
func TestAuthFlow(t *testing.T) {
	t.Log("Testing Signup -> OTP -> Login flow...")
	// TODO: Phase 2 - Spin up testcontainer with Postgres, apply migrations, run HTTP tests against router
	// This ensures the DB queries and JWT generation work harmoniously.
}
