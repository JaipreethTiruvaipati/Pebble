// Package queries encapsulates parameterized SQL access for Pebble domain tables.
// Penalty queries create consent-window penalties and move confirmed cash into the investment pool.
package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreatePendingPenalty inserts a penalties row awaiting user consent or auto-confirm.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: penalised user
//   - lineItemID: impulse line_items.id that triggered the penalty
//   - amount: INR to debit when confirmed
//   - consentHours: users.consent_hours; sets expires_at = now + consentHours
//
// Returns:
//   - uuid.UUID: new penalties.id
//   - error: insert failure
//
// Pebble flow: penalty-service after bills.scored; publishes wallet.penalty_queued until
// ConfirmExpiredPenalties or user action moves status to confirmed.
func CreatePendingPenalty(ctx context.Context, pool *pgxpool.Pool, userID, lineItemID uuid.UUID, amount float64, consentHours int) (uuid.UUID, error) {
	var id uuid.UUID
	expires := time.Now().Add(time.Duration(consentHours) * time.Hour)
	err := pool.QueryRow(ctx, `
		INSERT INTO penalties (user_id, line_item_id, amount, status, expires_at)
		VALUES ($1, $2, $3, 'pending', $4)
		RETURNING id`,
		userID, lineItemID, amount, expires,
	).Scan(&id)
	return id, err
}

// AddPoolContribution records confirmed penalty cash in pool_contributions for batch investment.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: contributor
//   - penaltyID: optional penalties.id link (nullable FK)
//   - amount: INR entering the pooled state
//
// Returns:
//   - error: insert failure
//
// Pebble flow: penalty-service after wallet debit on confirm; investment-service later
// picks up status='pooled' rows via SumPooledAmount / MarkPoolInvested.
func AddPoolContribution(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, penaltyID *uuid.UUID, amount float64) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO pool_contributions (user_id, penalty_id, amount, status)
		VALUES ($1, $2, $3, 'pooled')`,
		userID, penaltyID, amount,
	)
	return err
}

// ConfirmExpiredPenalties bulk-updates pending penalties past their consent window.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//
// Returns:
//   - int64: number of rows updated to status='confirmed'
//   - error: update failure
//
// Pebble flow: penalty-service cron; confirmed rows trigger wallet debits, pool contributions,
// and queue.TopicWalletPenaltyConfirmed events for downstream wallet-service.
func ConfirmExpiredPenalties(ctx context.Context, pool *pgxpool.Pool) (int64, error) {
	tag, err := pool.Exec(ctx, `
		UPDATE penalties
		SET status = 'confirmed', confirmed_at = NOW()
		WHERE status = 'pending' AND expires_at < NOW()`)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
