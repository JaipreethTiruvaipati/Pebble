package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/jaipreeth/pebble/backend/internal/models"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

// GeminiClient handles interactions with Google's Gemini API.
type GeminiClient struct {
	client *genai.Client
}

// ReceiptExtraction represents the structured data we expect back from Gemini.
type ReceiptExtraction struct {
	Merchant    string              `json:"merchant"`
	TotalAmount float64             `json:"total_amount"`
	Items       []models.ScoredItem `json:"items"`
}

// NewGeminiClient initializes a new GenAI client.
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

// Close closes the underlying Gemini client.
func (g *GeminiClient) Close() {
	g.client.Close()
}

// ExtractAndScore takes the raw messy text from Google Vision OCR and uses Gemini to extract 
// structured data and assign impulse scores in a single pass.
func (g *GeminiClient) ExtractAndScore(ctx context.Context, rawOCRText string) (*ReceiptExtraction, error) {
	// We use gemini-1.5-flash as it is extremely fast and great at JSON extraction
	model := g.client.GenerativeModel("gemini-1.5-flash")
	
	// Force the model to return valid JSON
	model.ResponseMIMEType = "application/json"
	
	prompt := `
You are an expert Indian financial assistant. You will be given raw, unstructured text extracted from an Indian shopping receipt via OCR.
Your job is to parse this text into a strict JSON format.

INSTRUCTIONS:
1. Identify the "merchant" name.
2. Identify the "total_amount" paid (numeric only).
3. Extract all individual line items. Do not include taxes (GST), service charges, or totals as items.
4. For each item, you must determine:
   - "name": the product name.
   - "amount": the cost of that specific item (numeric only).
   - "category": assign one of [food, beverage, essential, subscription, entertainment, transport, other].
   - "score": an Impulse Score from 0 to 100 (0 = highly essential, 100 = total impulse/luxury buy).
   - "reasoning": a very short 1-sentence explanation of why you gave this score.

RAW OCR TEXT:
` + rawOCRText

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gemini API call failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini returned an empty response")
	}

	// Extract the JSON text from the response
	var jsonStr string
	if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		jsonStr = string(text)
	} else {
		return nil, fmt.Errorf("unexpected response type from gemini")
	}

	// Clean up any markdown formatting (e.g. ` + "```" + `json ... ` + "```" + `) just in case
	jsonStr = strings.TrimPrefix(jsonStr, "```json\n")
	jsonStr = strings.TrimSuffix(jsonStr, "\n```")

	// Parse the JSON directly into our Go struct
	var result ReceiptExtraction
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Error().Err(err).Str("json", jsonStr).Msg("failed to parse gemini JSON")
		return nil, fmt.Errorf("failed to unmarshal gemini response: %w", err)
	}

	log.Info().Str("merchant", result.Merchant).Int("items", len(result.Items)).Msg("successfully parsed and scored receipt")
	return &result, nil
}
