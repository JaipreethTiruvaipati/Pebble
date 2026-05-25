package llm

import "testing"

func TestClampScore(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{-10, 0},
		{0, 0},
		{50, 50},
		{100, 100},
		{150, 100},
	}
	for _, tc := range tests {
		got := ClampScore(tc.input)
		if got != tc.want {
			t.Errorf("ClampScore(%d) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestFallbackScore(t *testing.T) {
	tests := []struct {
		category string
		want     int
	}{
		{"food", 35},
		{"essential", 10},
		{"entertainment", 65},
		{"unknown_category", 50},
	}
	for _, tc := range tests {
		got := FallbackScore(tc.category)
		if got != tc.want {
			t.Errorf("FallbackScore(%q) = %d, want %d", tc.category, got, tc.want)
		}
	}
}

func TestAdjustScoreForTime(t *testing.T) {
	tests := []struct {
		score int
		hour  int
		want  int
	}{
		{50, 14, 50},   // afternoon — no adjustment
		{50, 22, 65},   // 10 PM — +15
		{50, 2, 65},    // 2 AM — +15
		{90, 23, 100},  // late night + high score → clamped at 100
		{30, 6, 30},    // 6 AM — no adjustment (boundary)
	}
	for _, tc := range tests {
		got := AdjustScoreForTime(tc.score, tc.hour)
		if got != tc.want {
			t.Errorf("AdjustScoreForTime(%d, %d) = %d, want %d", tc.score, tc.hour, got, tc.want)
		}
	}
}
