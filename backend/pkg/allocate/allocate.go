package allocate

import "github.com/jaipreeth/pebble/backend/internal/models"

// AllocationResult represents the final percentage distribution of pooled funds.
type AllocationResult struct {
	Equity float64
	Gold   float64
	Bonds  float64
}

// ComputeAllocation calculates the optimal spread based on market signals.
// It applies clamping (min/max bounds) and normalizes the final array to equal 100%.
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
