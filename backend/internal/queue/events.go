package queue

import "github.com/google/uuid"

// Event topics
const (
	TopicBillsUploaded         = "bills.uploaded"
	TopicBillsScored           = "bills.scored"
	TopicWalletPenaltyQueued   = "wallet.penalty_queued"
	TopicWalletPenaltyConfirmed = "wallet.penalties_confirmed"
)

// PenaltyQueuedEvent is fired when a bill is scored and penalties are pending consent.
type PenaltyQueuedEvent struct {
	UserID        uuid.UUID `json:"user_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	TotalPending  float64   `json:"total_pending"`
}

// PenaltiesConfirmedEvent is fired when the consent window expires or the user manually confirms.
type PenaltiesConfirmedEvent struct {
	UserID      uuid.UUID   `json:"user_id"`
	PenaltyIDs  []uuid.UUID `json:"penalty_ids"`
	TotalAmount float64     `json:"total_amount"`
}
