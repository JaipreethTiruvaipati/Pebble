// Package queries encapsulates parameterized SQL access for Pebble domain tables.
// Transaction queries support bill logging, LLM scoring persistence, and penalty calculation.
package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/models"
)

// LineItemRow is a minimal line_items projection used by penalty-service arithmetic.
type LineItemRow struct {
	ID           uuid.UUID
	Amount       float64
	ImpulseScore int
}

// ListLineItemsForTransaction returns scored line items for a parent transaction.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - transactionID: line_items.transaction_id foreign key
//
// Returns:
//   - []LineItemRow: id, amount, impulse_score per row
//   - error: query or scan failure
//
// Pebble flow: penalty-service loads items after bills.scored to compute per-line penalties.
func ListLineItemsForTransaction(ctx context.Context, pool *pgxpool.Pool, transactionID uuid.UUID) ([]LineItemRow, error) {
	rows, err := pool.Query(ctx, `
		SELECT id, amount, impulse_score
		FROM line_items
		WHERE transaction_id = $1`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LineItemRow
	for rows.Next() {
		var li LineItemRow
		if err := rows.Scan(&li.ID, &li.Amount, &li.ImpulseScore); err != nil {
			return nil, err
		}
		items = append(items, li)
	}
	return items, rows.Err()
}

// GetUserPenaltySettings loads penalty_rate, penalty_threshold, and consent_hours for a user.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: users.id
//
// Returns:
//   - rate: base penalty fraction (0.05–0.30) before streak/referral discounts
//   - threshold: impulse_score cutoff for penalised line items
//   - consentHours: hours until pending penalties auto-confirm
//   - err: query failure
func GetUserPenaltySettings(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (rate float64, threshold, consentHours int, err error) {
	err = pool.QueryRow(ctx, `
		SELECT penalty_rate, penalty_threshold, consent_hours
		FROM users WHERE id = $1`, userID,
	).Scan(&rate, &threshold, &consentHours)
	return
}

// CreateTransaction inserts a manual UPI transaction in pending status.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: transaction owner
//   - merchant: payee name
//   - total: transaction total_amount in INR
//
// Returns:
//   - uuid.UUID: new transactions.id
//   - error: insert failure
//
// Pebble flow: bill-service after receipt upload or manual log; status='pending' until
// scoring completes and MarkTransactionScored runs.
func CreateTransaction(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, merchant string, total float64) (uuid.UUID, error) {
	var id uuid.UUID
	err := pool.QueryRow(ctx, `
		INSERT INTO transactions (user_id, merchant, total_amount, logged_at, status)
		VALUES ($1, $2, $3, NOW(), 'pending')
		RETURNING id`, userID, merchant, total,
	).Scan(&id)
	return id, err
}

// InsertLineItem persists one LLM-scored line item for a transaction.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - txID: parent transactions.id
//   - item: ScoredItem from Gemini (name, amount, score, category, reasoning)
//
// Returns:
//   - uuid.UUID: new line_items.id (referenced by penalties.line_item_id)
//   - error: insert failure
func InsertLineItem(ctx context.Context, pool *pgxpool.Pool, txID uuid.UUID, item models.ScoredItem) (uuid.UUID, error) {
	var id uuid.UUID
	err := pool.QueryRow(ctx, `
		INSERT INTO line_items (transaction_id, name, amount, impulse_score, category, reasoning)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		txID, item.Name, item.Amount, item.Score, item.Category, item.Reasoning,
	).Scan(&id)
	return id, err
}

// MarkTransactionScored sets scored_at and moves status to 'scored'.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - txID: transactions.id to update
//
// Returns:
//   - error: update failure
//
// Pebble flow: bill-service after all InsertLineItem calls; triggers penalty-service via
// queue.TopicBillsScored / PenaltyQueuedEvent.
func MarkTransactionScored(ctx context.Context, pool *pgxpool.Pool, txID uuid.UUID) error {
	_, err := pool.Exec(ctx, `
		UPDATE transactions SET scored_at = NOW(), status = 'scored' WHERE id = $1`, txID)
	return err
}
