//go:build integration

package integration

import (
	"testing"
)

func TestAuthFlow(t *testing.T) {
	t.Log("Testing Signup -> OTP -> Login flow...")
	// TODO: Phase 2 - Spin up testcontainer with Postgres, apply migrations, run HTTP tests against router
	// This ensures the DB queries and JWT generation work harmoniously.
}
