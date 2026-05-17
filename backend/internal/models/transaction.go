package models

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a single logged UPI payment.
// Maps to the 'transactions' table.
type Transaction struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Merchant    string     `json:"merchant"`
	TotalAmount float64    `json:"total_amount"`
	BillS3Key   *string    `json:"bill_s3_key,omitempty"` // Pointer because it can be null
	ScoredAt    *time.Time `json:"scored_at,omitempty"`   // Null until LLM is done
	LoggedAt    time.Time  `json:"logged_at"`
	
	// Relationships (often populated via joins/extra queries)
	LineItems []LineItem `json:"line_items,omitempty"`
}

// LineItem represents a single item extracted from a receipt and scored.
// Maps to the 'line_items' table.
type LineItem struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	Name          string    `json:"name"`
	Amount        float64   `json:"amount"`
	Quantity      int       `json:"quantity"`
	ImpulseScore  int       `json:"impulse_score"` // 0-100 assigned by LLM
	Category      string    `json:"category"`      // food, beverage, essential, etc.
	Reasoning     string    `json:"reasoning"`     // 1-sentence LLM justification
	
	// User override fields
	UserOverridden bool `json:"user_overridden"`
	OverrideScore  *int `json:"override_score,omitempty"` // User's manual score correction
}

// ScoredItem is a specialized struct used during the LLM scoring pipeline.
// It's the JSON structure we expect the Gemini API to return.
type ScoredItem struct {
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	Score     int     `json:"score"`
	Category  string  `json:"category"`
	Reasoning string  `json:"reasoning"`
}
