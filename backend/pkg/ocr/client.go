// Package ocr wraps cloud optical character recognition for Pebble receipt uploads.
//
// bill-service stores receipt images; scoring-service consumes BillUploaded events,
// runs OCR here, then passes text to pkg/llm Gemini for structuring and impulse scoring.
package ocr

import (
	"context"

	"github.com/rs/zerolog/log"
)

// VisionClient holds configuration for Google Cloud Vision API access.
// credPath is the service-account JSON path from config (GOOGLE_VISION_CRED_PATH).
type VisionClient struct {
	credPath string
}

// NewVisionClient constructs a VisionClient that will use credPath for API authentication.
// scoring-service instantiates one at startup alongside the Gemini client.
func NewVisionClient(credPath string) *VisionClient {
	return &VisionClient{credPath: credPath}
}

// ExtractText runs document text detection on imageBytes and returns plain text.
//
// Phase 1 stub returns fixed sample text for local development; production will call
// Vision API with c.credPath credentials. processBillUploaded passes the result to
// llm.GeminiClient.ExtractAndScore.
func (c *VisionClient) ExtractText(ctx context.Context, imageBytes []byte) (string, error) {
	log.Info().Int("bytes", len(imageBytes)).Msg("sending image to Google Vision API")
	// TODO: Phase 1 - Implement actual Google Cloud Vision API call
	// stub response
	return "Domino's Pizza\nTotal: Rs 850.00\nMargherita Pizza - 400\nCheese Burst - 450", nil
}
