// Package allocate implements signal-driven portfolio allocation for Pebble pooled investments.
//
// Penalty cash pooled by penalty-service is invested by investment-service. This package
// converts market signals (cached by market-poller in Redis) into percentage weights
// across equity, gold, and bonds before orders are sent to the broker layer.
package allocate

import "github.com/jaipreeth/pebble/backend/internal/models"

// AllocationResult holds the final percentage distribution of pooled funds across
// Pebble's three asset sleeves. Values are normalized to sum to 100 after clamping.
type AllocationResult struct {
	Equity float64
	Gold   float64
	Bonds  float64
}

// ComputeAllocation derives an equity/gold/bonds mix from market signals.
//
// It starts from a balanced baseline (40/20/40), applies per-signal adjustments
// (e.g. equity BUY shifts weight from bonds), enforces safety clamps (equity 10–70%,
// bonds and gold floors), then renormalizes so the three percentages sum to exactly 100.
//
// Called by ComputeBrokerOrders and indirectly by investment-service PoolExecutor
// when Redis has no cached signals and defaults are used upstream.
func ComputeAllocation(signals []models.MarketSignal) AllocationResult {
	// Base safe allocation (balanced)
	alloc := AllocationResult{
		Equity: 40.0,
		Gold:   20.0,
		Bonds:  40.0,
	}

	// 1. Apply signal adjustments
	for _, sig := range signals {
		if sig.AssetClass == "equity" && sig.Action == "BUY" {
			alloc.Equity += 15.0
			alloc.Bonds -= 15.0
		} else if sig.AssetClass == "gold" && sig.Action == "SELL" {
			alloc.Gold -= 10.0
			alloc.Bonds += 10.0
		}
	}

	// 2. Clamp constraints (Safety rails)
	// Don't let AI put >70% in Equity or <10% in Bonds
	if alloc.Equity > 70.0 {
		alloc.Equity = 70.0
	}
	if alloc.Equity < 10.0 {
		alloc.Equity = 10.0
	}
	if alloc.Bonds < 10.0 {
		alloc.Bonds = 10.0
	}
	if alloc.Gold < 5.0 {
		alloc.Gold = 5.0
	}

	// 3. Normalize to exactly 100.0%
	total := alloc.Equity + alloc.Gold + alloc.Bonds
	alloc.Equity = (alloc.Equity / total) * 100.0
	alloc.Gold = (alloc.Gold / total) * 100.0
	alloc.Bonds = (alloc.Bonds / total) * 100.0

	return alloc
}
