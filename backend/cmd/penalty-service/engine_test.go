package main

import (
	"math"
	"testing"
)

func TestEffectivePenaltyRate_StreakDiscount(t *testing.T) {
	tests := []struct {
		base     float64
		streak   int
		referral bool
		want     float64
	}{
		{0.15, 0, false, 0.15},
		{0.15, 4, false, 0.14},
		{0.15, 8, false, 0.13},
		{0.15, 0, true, 0.13},
		{0.06, 8, true, 0.05}, // floor
	}
	for _, tc := range tests {
		got := EffectivePenaltyRate(tc.base, tc.streak, tc.referral)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("EffectivePenaltyRate(%v, %d) = %v, want %v", tc.base, tc.streak, got, tc.want)
		}
	}
}
