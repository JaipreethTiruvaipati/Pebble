package models

import (
	"time"

	"github.com/google/uuid"
)

// Penalty represents a deduction to be made from a user's wallet due to an impulse buy.
type Penalty struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	LineItemID  uuid.UUID  `json:"line_item_id"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"` // pending, confirmed, contested, expired
	ExpiresAt   time.Time  `json:"expires_at"` // 24h consent window
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	ContestedAt *time.Time `json:"contested_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
