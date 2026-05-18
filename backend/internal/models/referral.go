// Package models defines domain structs that map to Pebble PostgreSQL tables and API JSON bodies.
// Referral types back the invite programme and referrer penalty discount.
package models

import (
	"time"

	"github.com/google/uuid"
)

// ReferralCode is a user's shareable invite code stored in referral_codes.
type ReferralCode struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

// ReferralStats summarises a referrer's programme status for GET /referrals/me.
type ReferralStats struct {
	Code            string  `json:"code"`
	RedemptionCount int     `json:"redemption_count"`
	DiscountPct     float64 `json:"discount_pct"` // active referrer discount (e.g. 2 when qualified)
}
