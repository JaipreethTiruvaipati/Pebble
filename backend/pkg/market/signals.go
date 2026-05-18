// Package market ingests Indian market data and produces allocation signals for Pebble.
//
// signals.go turns normalized market snapshots into models.MarketSignal slices consumed
// by pkg/allocate and cached for investment-service.
package market

import (
	"time"

	"github.com/jaipreeth/pebble/backend/internal/models"
)

// ComputeOpportunitySignals derives actionable signals from NSE, MCX, and CCIL snapshots.
//
// Each signal includes asset class, indicator name, numeric value, and recommended action
// (BUY, SELL, HOLD). market-poller stores the returned slice at cache.KeyMarketSignals;
// PoolExecutor reads it when executing pooled investments. Phase 2 will replace the stub
// with real technical analysis; nil maps are accepted today.
func ComputeOpportunitySignals(nseData, mcxData, ccilData map[string]interface{}) ([]models.MarketSignal, error) {
	// Stub implementation - Phase 2 will implement real technical analysis (RSI, MACD)
	signals := []models.MarketSignal{
		{
			AssetClass: "equity",
			Indicator:  "RSI",
			Value:      28.5, // Oversold
			Action:     "BUY",
			Timestamp:  time.Now(),
		},
		{
			AssetClass: "gold",
			Indicator:  "MACD",
			Value:      1.2,
			Action:     "HOLD",
			Timestamp:  time.Now(),
		},
		{
			AssetClass: "bonds",
			Indicator:  "YieldCurve",
			Value:      7.1,
			Action:     "BUY",
			Timestamp:  time.Now(),
		},
	}
	return signals, nil
}
