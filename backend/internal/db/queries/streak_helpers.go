// Package queries encapsulates parameterized SQL access for Pebble domain tables.
// Streak helpers compute effective penalty rates shared by penalty-service and api-gateway.
package queries

// EffectivePenaltyRateForUser applies streak and referrer discounts to the user's base penalty rate.
//
// Parameters:
//   - baseRate: users.penalty_rate before discounts (e.g. 0.15)
//   - streakCount: users.streak_count; every 4 weeks removes 1% (0.01) from rate
//   - hasReferrerDiscount: true when the user has ≥1 referral redemption as referrer
//
// Returns:
//   - float64: discounted rate, floored at 5% (0.05)
//
// How it works: discount = (streakCount/4)×0.01 + ReferrerDiscountPct when applicable;
// mirrors penalty-service engine logic so API previews match actual debits. Pure function—no DB I/O.
func EffectivePenaltyRateForUser(baseRate float64, streakCount int, hasReferrerDiscount bool) float64 {
	discount := float64(streakCount/4) * 0.01
	if hasReferrerDiscount {
		discount += ReferrerDiscountPct
	}
	rate := baseRate - discount
	if rate < 0.05 {
		return 0.05
	}
	return rate
}
