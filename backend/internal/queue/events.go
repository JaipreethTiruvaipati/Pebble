// Package queue defines RabbitMQ connectivity, event payloads, and publish/consume helpers
// used to decouple Pebble microservices (bill scoring → penalties → pool investment).
package queue

import "github.com/google/uuid"

// RabbitMQ routing keys published on the pebble.events topic exchange.
const (
	// TopicBillsUploaded is emitted when a receipt image lands in S3 and bill-service should score it.
	TopicBillsUploaded = "bills.uploaded"
	// TopicBillsScored is emitted when LLM scoring finishes; penalty-service creates pending penalties.
	TopicBillsScored = "bills.scored"
	// TopicWalletPenaltyQueued is emitted when pending penalties are created for a transaction.
	TopicWalletPenaltyQueued = "wallet.penalty_queued"
	// TopicWalletPenaltyConfirmed is emitted when penalties confirm and wallet debits should run.
	TopicWalletPenaltyConfirmed = "wallet.penalties_confirmed"
	// TopicInvestmentsExecuted is emitted after MarkPoolInvested completes a broker batch.
	TopicInvestmentsExecuted = "investments.executed"
	// TopicStreakUpdated is emitted when scoring-service awards a low-impulse week streak.
	TopicStreakUpdated = "streak.updated"
)

// PenaltyQueuedEvent notifies wallet/notification services that consent-window penalties exist.
type PenaltyQueuedEvent struct {
	UserID        uuid.UUID `json:"user_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	TotalPending  float64   `json:"total_pending"`
}

// PenaltiesConfirmedEvent signals that penalties moved to confirmed and should debit wallets.
type PenaltiesConfirmedEvent struct {
	UserID      uuid.UUID   `json:"user_id"`
	PenaltyIDs  []uuid.UUID `json:"penalty_ids"`
	TotalAmount float64     `json:"total_amount"`
}

// InvestmentAllocation describes one asset-class slice of a completed pool investment batch.
type InvestmentAllocation struct {
	AssetClass string  `json:"asset_class"`
	Amount     float64 `json:"amount"`
}

// BillsScoredEvent triggers penalty calculation after line items are persisted.
type BillsScoredEvent struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	UserID        uuid.UUID `json:"user_id"`
}

// StreakUpdatedEvent notifies clients when a user earns another low-impulse week.
type StreakUpdatedEvent struct {
	UserID       uuid.UUID `json:"user_id"`
	StreakCount  int       `json:"streak_count"`
	WeekAvgScore float64   `json:"week_avg_score"`
}

// InvestmentExecutedEvent summarizes a completed ExecutePool / MarkPoolInvested run.
type InvestmentExecutedEvent struct {
	TriggerType   string                 `json:"trigger_type"`
	TotalAmount   float64                `json:"total_amount"`
	BrokerRef     string                 `json:"broker_ref"`
	InvestmentIDs []uuid.UUID            `json:"investment_ids"`
	Allocation    []InvestmentAllocation `json:"allocation"`
}
