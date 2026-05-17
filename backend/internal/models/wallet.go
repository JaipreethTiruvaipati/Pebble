package models

import (
	"time"

	"github.com/google/uuid"
)

// WalletTransaction records all money movements in the user's wallet.
type WalletTransaction struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	Type         string     `json:"type"` // topup, penalty_debit, investment_debit, refund
	Amount       float64    `json:"amount"`
	ReferenceID  *uuid.UUID `json:"reference_id,omitempty"` // ID of penalty or investment
	BalanceAfter float64    `json:"balance_after"`
	CreatedAt    time.Time  `json:"created_at"`
}
