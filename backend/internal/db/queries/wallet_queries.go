package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WalletSnapshot is the wallet card data for GET /wallet/balance.
type WalletSnapshot struct {
	Balance       float64 `json:"balance"`
	PendingTotal  float64 `json:"pending_total"`
	InvestedTotal float64 `json:"invested_total"`
}

// WalletLedgerEntry is one row in GET /wallet/ledger.
type WalletLedgerEntry struct {
	ID           uuid.UUID `json:"id"`
	Type         string    `json:"type"`
	Amount       float64   `json:"amount"`
	BalanceAfter float64   `json:"balance_after"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetWalletSnapshot loads wallet balances and sums pending penalties from DB.
func GetWalletSnapshot(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*WalletSnapshot, error) {
	var snap WalletSnapshot
	err := pool.QueryRow(ctx, `
		SELECT COALESCE(balance, 0), COALESCE(invested_total, 0)
		FROM wallets WHERE user_id = $1`, userID,
	).Scan(&snap.Balance, &snap.InvestedTotal)
	if err == pgx.ErrNoRows {
		return &WalletSnapshot{}, nil
	}
	if err != nil {
		return nil, err
	}
	_ = pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0) FROM penalties
		WHERE user_id = $1 AND status = 'pending'`, userID,
	).Scan(&snap.PendingTotal)
	return &snap, nil
}

// ListWalletLedger returns recent wallet_transactions for a user.
func ListWalletLedger(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, limit int) ([]WalletLedgerEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := pool.Query(ctx, `
		SELECT id, type, amount, balance_after, created_at
		FROM wallet_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []WalletLedgerEntry
	for rows.Next() {
		var e WalletLedgerEntry
		if err := rows.Scan(&e.ID, &e.Type, &e.Amount, &e.BalanceAfter, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// CreditWalletTopup adds funds for dev Razorpay top-up simulation.
func CreditWalletTopup(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, amount float64) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var balance float64
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(balance, 0) FROM wallets WHERE user_id = $1 FOR UPDATE`, userID,
	).Scan(&balance)
	if err == pgx.ErrNoRows {
		_, err = tx.Exec(ctx, `INSERT INTO wallets (user_id, balance, topup_total) VALUES ($1, $2, $2)`, userID, amount)
		if err != nil {
			return err
		}
		balance = 0
	} else if err != nil {
		return err
	}

	newBal := balance + amount
	_, err = tx.Exec(ctx, `
		UPDATE wallets SET balance = $2, topup_total = topup_total + $3, updated_at = NOW()
		WHERE user_id = $1`, userID, newBal, amount)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions (user_id, type, amount, balance_after)
		VALUES ($1, 'topup', $2, $3)`, userID, amount, newBal)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
