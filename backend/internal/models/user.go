package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user in the system.
// Maps to the 'users' table in the database.
type User struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	PasswordHash     string    `json:"-"` // Never expose in JSON responses
	RiskProfile      string    `json:"risk_profile"` // conservative, moderate, aggressive
	PenaltyRate      float64   `json:"penalty_rate"` // 0.05 - 0.30 (5% to 30%)
	PenaltyThreshold int       `json:"penalty_threshold"` // 30 - 90
	InvestThreshold  float64   `json:"invest_threshold"` // amount required to trigger investment
	ConsentHours     int       `json:"consent_hours"` // how long user has to contest penalty
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RefreshToken represents a JWT refresh token session.
// Depending on implementation, you might store these in DB or Redis.
type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"` // Hashed token, just like passwords
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
