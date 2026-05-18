// Package cache wraps Redis for Pebble: connection lifecycle, JSON get/set, and
// key helpers. This file centralizes Redis key names and TTLs so market-poller,
// api-gateway, and investment paths stay consistent and avoid key collisions.
package cache

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Market and portfolio Redis key names and TTLs used across Pebble services.
const (
	// KeyMarketNSE is the Redis key for cached NSE index levels written by market-poller.
	KeyMarketNSE = "market:nse:indices"
	// KeyMarketMCX is the Redis key for cached MCX gold quotes.
	KeyMarketMCX = "market:mcx:gold"
	// KeyMarketCCIL is the Redis key for cached CCIL government bond yields.
	KeyMarketCCIL = "market:ccil:yields"
	// KeyMarketAMFI is the Redis key for cached AMFI mutual fund NAVs.
	KeyMarketAMFI = "market:amfi:navs"
	// KeyMarketSignals is the Redis key for the latest aggregated market signals payload.
	KeyMarketSignals = "market:signals:latest"

	// PortfolioSummaryTTL is how long a user's portfolio summary remains in Redis before expiry.
	// Chosen to keep api-gateway p95 under ~200ms on hot reads while staying fresh enough after trades.
	PortfolioSummaryTTL = 30 * time.Second
)

// KeyPortfolioSummary returns the per-user Redis key for a cached portfolio digest.
//
// Inputs: userID from auth context or database row.
// Outputs: key "portfolio:summary:<uuid>" for SetJSON/GetJSON/Delete in api-gateway.
//
// Invalidated on investment or bill events that change holdings or allocation views.
func KeyPortfolioSummary(userID uuid.UUID) string {
	return fmt.Sprintf("portfolio:summary:%s", userID)
}
