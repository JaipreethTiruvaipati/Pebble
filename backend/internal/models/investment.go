// Package models defines domain structs that map to Pebble PostgreSQL tables and API JSON bodies.
// Investment types represent pooled penalty cash deployed into market assets.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Investment is a broker purchase on behalf of a user after pool batch execution.
// Maps to investments; created by queries.MarkPoolInvested from pooled contributions.
type Investment struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AssetClass    string    `json:"asset_class"` // e.g., 'equity', 'gold', 'bonds'
	Amount        float64   `json:"amount"`
	Units         float64   `json:"units"`
	NAVAtPurchase float64   `json:"nav_at_purchase"`
	Status        string    `json:"status"`
	TriggerType   string    `json:"trigger_type,omitempty"`
	BrokerRef     string    `json:"broker_ref,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// PoolContribution links confirmed penalty cash to the shared investment pool.
// Maps to pool_contributions; status moves pooled → invested after MarkPoolInvested.
type PoolContribution struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	PenaltyID  *uuid.UUID `json:"penalty_id,omitempty"`
	Amount     float64    `json:"amount"`
	Status     string     `json:"status"` // 'pooled', 'invested'
	InvestedAt *time.Time `json:"invested_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// MarketSignal is an AI/market-poller decision that can trigger discretionary pool invests.
// Used by market-poller and investment-service allocation logic.
type MarketSignal struct {
	AssetClass string    `json:"asset_class"`
	Indicator  string    `json:"indicator"` // e.g., "RSI", "YieldCurve"
	Value      float64   `json:"value"`
	Action     string    `json:"action"` // "BUY", "SELL", "HOLD"
	Timestamp  time.Time `json:"timestamp"`
}
