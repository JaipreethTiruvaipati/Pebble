package models

import (
	"time"

	"github.com/google/uuid"
)

// Investment represents an asset purchased by Pebble on behalf of the user pool.
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

// PoolContribution represents a confirmed penalty that has been added to the investment pool.
type PoolContribution struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	PenaltyID  *uuid.UUID `json:"penalty_id,omitempty"`
	Amount     float64    `json:"amount"`
	Status     string     `json:"status"` // 'pooled', 'invested'
	InvestedAt *time.Time `json:"invested_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// MarketSignal represents an AI-generated decision based on live market data.
type MarketSignal struct {
	AssetClass string    `json:"asset_class"`
	Indicator  string    `json:"indicator"` // e.g., "RSI", "YieldCurve"
	Value      float64   `json:"value"`
	Action     string    `json:"action"` // "BUY", "SELL", "HOLD"
	Timestamp  time.Time `json:"timestamp"`
}
