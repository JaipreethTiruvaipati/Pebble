// Package llm provides large-language-model integrations for Pebble receipt intelligence.
//
// scorer.go hosts pure scoring helpers (category heuristics, score clamping) that complement
// Gemini extraction or serve as fallback when the API is unavailable.
package llm

// CategoryHeuristic maps common merchant/item keywords to baseline impulse scores.
// Used as a fallback when Gemini is unavailable or for pre-screening before LLM calls.
var CategoryHeuristic = map[string]int{
	"food":          35,
	"beverage":      40,
	"essential":     10,
	"subscription":  50,
	"entertainment": 65,
	"transport":     15,
	"other":         50,
}

// ClampScore ensures an impulse score stays within the valid 0–100 range.
func ClampScore(score int) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

// FallbackScore returns a heuristic impulse score for a category when the LLM is unavailable.
// Falls back to 50 (neutral) for unknown categories.
func FallbackScore(category string) int {
	if score, ok := CategoryHeuristic[category]; ok {
		return score
	}
	return 50
}

// AdjustScoreForTime adds a penalty for late-night purchases (10 PM – 6 AM).
// Purchases during these hours are typically more impulsive.
func AdjustScoreForTime(score int, hour int) int {
	if hour >= 22 || hour < 6 {
		score += 15
	}
	return ClampScore(score)
}
