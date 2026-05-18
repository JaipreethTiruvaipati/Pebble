package queue

import "github.com/google/uuid"

// Event topics (RabbitMQ routing keys on pebble.events exchange)
const (
	TopicBillsUploaded          = "bills.uploaded"
	TopicBillsScored            = "bills.scored"
	TopicWalletPenaltyQueued    = "wallet.penalty_queued"
	TopicWalletPenaltyConfirmed = "wallet.penalties_confirmed"
	TopicInvestmentsExecuted    = "investments.executed"
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

// InvestmentAllocation is one sleeve of a completed investment batch.
type InvestmentAllocation struct {
	AssetClass string  `json:"asset_class"`
	Amount     float64 `json:"amount"`
}

// InvestmentExecutedEvent is fired after ExecutePool completes broker orders.
type InvestmentExecutedEvent struct {
	TriggerType   string                 `json:"trigger_type"`
	TotalAmount   float64                `json:"total_amount"`
	BrokerRef     string                 `json:"broker_ref"`
	InvestmentIDs []uuid.UUID            `json:"investment_ids"`
	Allocation    []InvestmentAllocation `json:"allocation"`
}
