// Package llm provides large-language-model integrations for Pebble receipt intelligence.
//
// prompts.go centralizes prompt templates for Gemini (and future models) so scoring-service
// can version and test extraction/scoring instructions independently of HTTP client code.
package llm

// ReceiptExtractionPrompt is the system prompt for extracting and scoring receipt line items.
// It instructs the LLM to return structured JSON with merchant, items, impulse scores (0–100),
// and reasoning for each item classification.
const ReceiptExtractionPrompt = `You are an expert Indian financial assistant. You will be given raw, unstructured text extracted from an Indian shopping receipt via OCR.
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

SCORING GUIDELINES:
- Essential groceries, medicine, rent, bills: 0–20
- Regular meals, commute, gym: 20–40
- Dining out, streaming subscriptions: 40–60
- Fashion, gadgets, impulsive online orders: 60–80
- Luxury items, late-night shopping, repeat unnecessary purchases: 80–100
- Consider time of purchase if visible (late night = higher score)
- Consider frequency context: "3rd pair this quarter" = higher score

RAW OCR TEXT:
`

// ScoreAdjustmentPrompt is used for re-scoring items when additional context is available
// (e.g., user purchase history, time of day patterns).
const ScoreAdjustmentPrompt = `You are a behavioral finance expert. Given a user's purchase history context and a new line item, adjust the impulse score (0-100) considering:
1. Frequency: Is this a repeat impulse purchase?
2. Timing: Late-night (10 PM–6 AM) purchases score higher.
3. Category saturation: Too many items in the same discretionary category this month.
4. Budget alignment: Does this fit the user's stated financial goals?

Return JSON: {"adjusted_score": <int>, "adjustment_reason": "<string>"}

CONTEXT:
`
