package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionSummary is a list row for GET /transactions.
type TransactionSummary struct {
	ID          uuid.UUID `json:"id"`
	Merchant    string    `json:"merchant"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	LoggedAt    time.Time `json:"logged_at"`
	AvgScore    float64   `json:"avg_score"`
	TotalPenalty float64  `json:"total_penalty"`
}

// LineItemDetail is a scored line item on a transaction.
type LineItemDetail struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Amount        float64   `json:"amount"`
	ImpulseScore  int       `json:"impulse_score"`
	Category      string    `json:"category"`
	Reasoning     string    `json:"reasoning"`
	UserOverridden bool     `json:"user_overridden"`
}

// TransactionDetail is GET /transactions/{id}.
type TransactionDetail struct {
	ID          uuid.UUID        `json:"id"`
	Merchant    string           `json:"merchant"`
	TotalAmount float64          `json:"total_amount"`
	Status      string           `json:"status"`
	LoggedAt    time.Time        `json:"logged_at"`
	LineItems   []LineItemDetail `json:"line_items"`
}

// ListTransactions returns recent transactions with aggregate scores.
func ListTransactions(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, limit int) ([]TransactionSummary, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := pool.Query(ctx, `
		SELECT t.id, t.merchant, t.total_amount, t.status, t.logged_at,
		       COALESCE(AVG(li.impulse_score), 0),
		       COALESCE(SUM(p.amount), 0)
		FROM transactions t
		LEFT JOIN line_items li ON li.transaction_id = t.id
		LEFT JOIN penalties p ON p.line_item_id = li.id
		WHERE t.user_id = $1
		GROUP BY t.id
		ORDER BY t.logged_at DESC
		LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TransactionSummary
	for rows.Next() {
		var t TransactionSummary
		if err := rows.Scan(&t.ID, &t.Merchant, &t.TotalAmount, &t.Status, &t.LoggedAt, &t.AvgScore, &t.TotalPenalty); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// GetTransactionDetail loads a transaction and line items scoped to the user.
func GetTransactionDetail(ctx context.Context, pool *pgxpool.Pool, userID, txID uuid.UUID) (*TransactionDetail, error) {
	var d TransactionDetail
	err := pool.QueryRow(ctx, `
		SELECT id, merchant, total_amount, status, logged_at
		FROM transactions WHERE id = $1 AND user_id = $2`, txID, userID,
	).Scan(&d.ID, &d.Merchant, &d.TotalAmount, &d.Status, &d.LoggedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rows, err := pool.Query(ctx, `
		SELECT id, name, amount, impulse_score, COALESCE(category,''), COALESCE(reasoning,''), user_overridden
		FROM line_items WHERE transaction_id = $1 ORDER BY created_at`, txID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var li LineItemDetail
		if err := rows.Scan(&li.ID, &li.Name, &li.Amount, &li.ImpulseScore, &li.Category, &li.Reasoning, &li.UserOverridden); err != nil {
			return nil, err
		}
		d.LineItems = append(d.LineItems, li)
	}
	return &d, rows.Err()
}

// OverrideLineItemScore sets user override on a line item owned via transaction.
func OverrideLineItemScore(ctx context.Context, pool *pgxpool.Pool, userID, lineItemID uuid.UUID, score int) error {
	tag, err := pool.Exec(ctx, `
		UPDATE line_items li
		SET user_overridden = TRUE, override_score = $3, impulse_score = $3
		FROM transactions t
		WHERE li.id = $1 AND li.transaction_id = t.id AND t.user_id = $2`,
		lineItemID, userID, score)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
