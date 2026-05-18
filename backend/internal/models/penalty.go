// Package models defines domain structs that map to Pebble PostgreSQL tables and API JSON bodies.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Penalty is an impulse-buy deduction with a consent window before wallet debit.
// Maps to penalties; confirmed amounts flow to pool_contributions for investment.
type Penalty struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	LineItemID  uuid.UUID  `json:"line_item_id"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"` // pending, confirmed, contested, expired
	ExpiresAt   time.Time  `json:"expires_at"` // consent window end
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	ContestedAt *time.Time `json:"contested_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
