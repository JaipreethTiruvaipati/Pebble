// Package models defines domain structs that map to Pebble PostgreSQL tables and API JSON bodies.
package models

import (
	"time"

	"github.com/google/uuid"
)

// User is a registered Pebble account with penalty settings and behaviour streak state.
// Maps to the users table; drives penalty rate, invest threshold, and cohort benchmarks.
type User struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	PasswordHash     string    `json:"-"` // Never expose in JSON responses
	RiskProfile      string    `json:"risk_profile"` // conservative, moderate, aggressive
	PenaltyRate      float64   `json:"penalty_rate"` // 0.05 - 0.30 (5% to 30%)
	PenaltyThreshold int       `json:"penalty_threshold"` // 30 - 90 impulse score cutoff
	InvestThreshold  float64   `json:"invest_threshold"` // pooled INR required to trigger batch invest
	ConsentHours     int       `json:"consent_hours"` // hours before pending penalties auto-confirm
	StreakCount      int       `json:"streak_count"`
	StreakLastUpdated *time.Time `json:"streak_last_updated,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RefreshToken represents a persisted JWT refresh session for auth rotation.
// May be stored in PostgreSQL or Redis depending on auth implementation.
type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"` // Hashed token, never returned to clients
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
