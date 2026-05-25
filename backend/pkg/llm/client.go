// Package llm provides large-language-model integrations for Pebble receipt intelligence.
//
// client.go defines a provider-agnostic Extractor interface so scoring-service can depend
// on abstractions rather than Gemini-specific types — enabling mock testing and future
// provider swaps (e.g. Claude, GPT-4o).
package llm

import "context"

// Extractor abstracts receipt parsing and impulse scoring from any LLM provider.
// scoring-service depends on this interface; GeminiClient implements it.
type Extractor interface {
	// ExtractAndScore parses raw OCR text into structured merchant data with impulse scores.
	ExtractAndScore(ctx context.Context, rawOCRText string) (*ReceiptExtraction, error)
	// Close releases provider SDK resources.
	Close()
}

// Ensure GeminiClient implements Extractor at compile time.
var _ Extractor = (*GeminiClient)(nil)
