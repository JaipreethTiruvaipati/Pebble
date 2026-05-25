// Package main (wallet.go) handles wallet ledger mutations triggered by confirmed penalties.
// When penalties expire the consent window, the confirmed amount is debited from the user's
// wallet balance and the funds are transferred to the investment pool.
package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// DebitWalletForPenalty atomically reduces a user's wallet balance by the penalty amount.
// Returns an error if the wallet has insufficient funds (the penalty is still recorded but
// the debit is deferred until the next top-up).
//
// In production, this would be wrapped in a database transaction with the pool contribution
// insert to ensure atomicity.
func DebitWalletForPenalty(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID, amount float64) error {
	tag, err := db.Exec(ctx, `
		UPDATE wallets 
		SET balance = balance - $1, 
		    updated_at = NOW() 
		WHERE user_id = $2 
		  AND balance >= $1`,
		amount, userID,
	)
	if err != nil {
		return fmt.Errorf("wallet debit failed: %w", err)
	}

	if tag.RowsAffected() == 0 {
		log.Warn().
			Str("user_id", userID.String()).
			Float64("amount", amount).
			Msg("insufficient wallet balance for penalty debit — deferred")
		return fmt.Errorf("insufficient wallet balance for user %s", userID)
	}

	log.Info().
		Str("user_id", userID.String()).
		Float64("amount", amount).
		Msg("wallet debited for confirmed penalty")
	return nil
}
