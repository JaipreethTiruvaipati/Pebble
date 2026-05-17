package ocr

import (
	"context"

	"github.com/rs/zerolog/log"
)

// VisionClient wraps the Google Cloud Vision API.
type VisionClient struct {
	credPath string
}

// NewVisionClient initializes the OCR client.
func NewVisionClient(credPath string) *VisionClient {
	return &VisionClient{credPath: credPath}
}

// ExtractText takes an image byte array and returns the raw text extracted by Google Vision.
func (c *VisionClient) ExtractText(ctx context.Context, imageBytes []byte) (string, error) {
	log.Info().Int("bytes", len(imageBytes)).Msg("sending image to Google Vision API")
	// TODO: Phase 1 - Implement actual Google Cloud Vision API call
	// stub response
	return "Domino's Pizza\nTotal: Rs 850.00\nMargherita Pizza - 400\nCheese Burst - 450", nil
}
