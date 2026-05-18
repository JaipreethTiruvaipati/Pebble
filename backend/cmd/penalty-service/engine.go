// Package main (engine.go) implements pure penalty math shared with api-gateway profile
// responses: effective rate discounts and per-line-item penalty amounts in INR.
package main

import "math"

// EffectivePenaltyRate applies streak (-1% per 4 weeks) and referrer (-2%) discounts with a 5% floor.
func EffectivePenaltyRate(baseRate float64, streakCount int, hasReferrerDiscount bool) float64 {
	discount := float64(streakCount/4) * 0.01
	if hasReferrerDiscount {
		discount += 0.02
	}
	rate := baseRate - discount
	if rate < 0.05 {
		return 0.05
	}
	return rate
}

// CalculatePenalty returns zero when impulseScore is below threshold; otherwise
// amount × penaltyRate × (score/100), rounded, capped at ₹500 and floored at ₹5.
func CalculatePenalty(itemAmount float64, impulseScore, threshold int, penaltyRate float64) float64 {
	if impulseScore < threshold {
		return 0
	}
	raw := itemAmount * penaltyRate * (float64(impulseScore) / 100.0)
	if raw < 5.0 {
		return 0
	}
	if raw > 500.0 {
		raw = 500.0
	}
	return math.Round(raw*100) / 100
}
