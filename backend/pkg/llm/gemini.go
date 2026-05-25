// Package llm provides large-language-model integrations for Pebble receipt intelligence.
//
// gemini.go implements the Google Gemini client used by scoring-service after OCR:
// raw receipt text in, structured merchant/items and impulse scores out.
package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/jaipreeth/pebble/backend/internal/models"
	"github.com/jaipreeth/pebble/backend/pkg/retry"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

// GeminiClient wraps the official Google Generative AI SDK for Pebble.
// scoring-service creates one at startup from GEMINI_API_KEY and passes it to processBillUploaded.
type GeminiClient struct {
	client *genai.Client
}

// ReceiptExtraction is the JSON shape Gemini returns for a single receipt parse pass.
// Items use models.ScoredItem so results map directly to queries.InsertLineItem in scoring-service.
type ReceiptExtraction struct {
	Merchant    string              `json:"merchant"`
	TotalAmount float64             `json:"total_amount"`
	Items       []models.ScoredItem `json:"items"`
}

// NewGeminiClient authenticates with apiKey and returns a ready GeminiClient.
//
// Returns an error if apiKey is empty or the underlying genai client cannot be created.
// scoring-service logs a warning and continues with OCR-only stubs when this fails.
func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &GeminiClient{client: client}, nil
}

// Close releases SDK resources. scoring-service defers this after a successful NewGeminiClient.
func (g *GeminiClient) Close() {
	g.client.Close()
}

// ExtractAndScore sends rawOCRText (from pkg/ocr VisionClient) to gemini-1.5-flash with a
// structured prompt and application/json response mode.
//
// The model extracts merchant, total, and line items with category, impulse score (0–100),
// and reasoning. Output is unmarshaled into ReceiptExtraction for persistence by scoring-service.
// Markdown code fences around JSON are stripped before parsing.
//
// Retries up to 3 times with exponential backoff on transient API failures.
func (g *GeminiClient) ExtractAndScore(ctx context.Context, rawOCRText string) (*ReceiptExtraction, error) {
	// We use gemini-1.5-flash as it is extremely fast and great at JSON extraction
	model := g.client.GenerativeModel("gemini-1.5-flash")

	// Force the model to return valid JSON
	model.ResponseMIMEType = "application/json"

	prompt := ReceiptExtractionPrompt + rawOCRText

	var result ReceiptExtraction

	err := retry.Do(ctx, retry.Default(), "gemini-extract-and-score", func() error {
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			return fmt.Errorf("gemini API call failed: %w", err)
		}

		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			return fmt.Errorf("gemini returned an empty response")
		}

		// Extract the JSON text from the response
		var jsonStr string
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			jsonStr = string(text)
		} else {
			return fmt.Errorf("unexpected response type from gemini")
		}

		// Clean up any markdown formatting (e.g. ```json ... ```) just in case
		jsonStr = strings.TrimPrefix(jsonStr, "```json\n")
		jsonStr = strings.TrimSuffix(jsonStr, "\n```")

		// Parse the JSON directly into our Go struct
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			log.Error().Err(err).Str("json", jsonStr).Msg("failed to parse gemini JSON")
			return fmt.Errorf("failed to unmarshal gemini response: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Clamp all scores to valid range
	for i := range result.Items {
		result.Items[i].Score = ClampScore(result.Items[i].Score)
	}

	log.Info().Str("merchant", result.Merchant).Int("items", len(result.Items)).Msg("successfully parsed and scored receipt")
	return &result, nil
}

