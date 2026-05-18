package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PenaltyRow is a penalty list item for the API.
type PenaltyRow struct {
	ID         uuid.UUID `json:"id"`
	LineItemID uuid.UUID `json:"line_item_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	ExpiresAt  time.Time `json:"expires_at"`
	Merchant   string    `json:"merchant,omitempty"`
	ItemName   string    `json:"item_name,omitempty"`
}

// PendingPenaltyBanner is the top pending penalty for the dashboard.
type PendingPenaltyBanner struct {
	Amount    float64   `json:"amount"`
	Source    string    `json:"source"`
	ExpiresAt time.Time `json:"expires_at"`
	PenaltyID uuid.UUID `json:"penalty_id"`
}

// ListPenalties returns penalties for a user filtered by status.
func ListPenalties(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, status string) ([]PenaltyRow, error) {
	q := `
		SELECT p.id, p.line_item_id, p.amount, p.status, p.expires_at,
		       COALESCE(t.merchant, ''), COALESCE(li.name, '')
		FROM penalties p
		JOIN line_items li ON li.id = p.line_item_id
		JOIN transactions t ON t.id = li.transaction_id
		WHERE p.user_id = $1`
	args := []interface{}{userID}
	if status != "" {
		q += ` AND p.status = $2`
		args = append(args, status)
	}
	q += ` ORDER BY p.created_at DESC LIMIT 100`

	rows, err := pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PenaltyRow
	for rows.Next() {
		var p PenaltyRow
		if err := rows.Scan(&p.ID, &p.LineItemID, &p.Amount, &p.Status, &p.ExpiresAt, &p.Merchant, &p.ItemName); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetTopPendingPenalty returns the largest pending penalty for dashboard banner.
func GetTopPendingPenalty(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*PendingPenaltyBanner, error) {
	var b PendingPenaltyBanner
	err := pool.QueryRow(ctx, `
		SELECT p.id, p.amount, p.expires_at,
		       COALESCE(li.name, t.merchant, 'Impulse spend')
		FROM penalties p
		JOIN line_items li ON li.id = p.line_item_id
		JOIN transactions t ON t.id = li.transaction_id
		WHERE p.user_id = $1 AND p.status = 'pending'
		ORDER BY p.amount DESC
		LIMIT 1`, userID,
	).Scan(&b.PenaltyID, &b.Amount, &b.ExpiresAt, &b.Source)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ContestPenalty marks a penalty as contested.
func ContestPenalty(ctx context.Context, pool *pgxpool.Pool, userID, penaltyID uuid.UUID) error {
	tag, err := pool.Exec(ctx, `
		UPDATE penalties SET status = 'contested', contested_at = NOW()
		WHERE id = $1 AND user_id = $2 AND status = 'pending'`, penaltyID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// ConfirmPenaltyEarly confirms a pending penalty before expiry.
func ConfirmPenaltyEarly(ctx context.Context, pool *pgxpool.Pool, userID, penaltyID uuid.UUID) error {
	tag, err := pool.Exec(ctx, `
		UPDATE penalties SET status = 'confirmed', confirmed_at = NOW()
		WHERE id = $1 AND user_id = $2 AND status = 'pending'`, penaltyID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
