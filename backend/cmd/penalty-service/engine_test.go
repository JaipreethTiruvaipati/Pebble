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

func TestCalculatePenalty_BelowThreshold(t *testing.T) {
	amt := CalculatePenalty(1000, 30, 50, 0.15)
	if amt != 0 {
		t.Fatalf("expected 0 penalty below threshold, got %v", amt)
	}
}

func TestCalculatePenalty_AboveThreshold(t *testing.T) {
	// 1000 * 0.15 * (72/100) = 108
	amt := CalculatePenalty(1000, 72, 50, 0.15)
	if amt != 108.0 {
		t.Fatalf("expected 108.0, got %v", amt)
	}
}

func TestCalculatePenalty_CappedAt500(t *testing.T) {
	// 50000 * 0.15 * (90/100) = 6750 → capped at 500
	amt := CalculatePenalty(50000, 90, 50, 0.15)
	if amt != 500.0 {
		t.Fatalf("expected 500.0 (cap), got %v", amt)
	}
}

func TestCalculatePenalty_FloorAt5(t *testing.T) {
	// 50 * 0.05 * (55/100) = 1.375 → below 5 floor → returns 0
	amt := CalculatePenalty(50, 55, 50, 0.05)
	if amt != 0 {
		t.Fatalf("expected 0 (below floor), got %v", amt)
	}
}

func TestCalculatePenalty_ExactThreshold(t *testing.T) {
	// Score exactly at threshold should NOT trigger penalty (< is strict)
	amt := CalculatePenalty(1000, 50, 50, 0.15)
	if amt != 0 {
		t.Fatalf("expected 0 at exact threshold, got %v", amt)
	}
}

func TestEffectivePenaltyRate_MaxStreakDiscount(t *testing.T) {
	// 52-week streak (13 * 1% = 13%) on 0.15 base → 0.02 → clamped to floor 0.05
	got := EffectivePenaltyRate(0.15, 52, false)
	if got != 0.05 {
		t.Fatalf("expected floor 0.05 for 52-week streak, got %v", got)
	}
}
