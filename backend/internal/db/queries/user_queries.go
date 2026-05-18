// Package queries encapsulates parameterized SQL access for Pebble domain tables.
// User queries support profile APIs, dev bootstrap, and scoring-service streak evaluation.
package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/models"
)

// EnsureDevUser inserts a placeholder users row so local JWT subjects satisfy foreign keys.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: UUID from dev JWT (sub claim)
//   - email: optional email; when empty, generates dev-{uuid-prefix}@pebble.in
//
// Returns:
//   - error: insert failure (ON CONFLICT DO NOTHING never errors on duplicate)
//
// Pebble flow: api-gateway dev auth middleware calls this before any user-scoped write
// so transactions, penalties, and wallets can reference users.id.
func EnsureDevUser(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, email string) error {
	if email == "" {
		email = fmt.Sprintf("dev-%s@pebble.in", userID.String()[:8])
	}
	phone := fmt.Sprintf("+91%09d", userID.ID()%1000000000)
	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, phone, password_hash)
		VALUES ($1, $2, $3, 'dev-hash')
		ON CONFLICT (id) DO NOTHING`,
		userID, email, phone,
	)
	return err
}

// EnsureDevWallet creates a wallets row for the user when missing (dev/local only).
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: wallet owner
//
// Returns:
//   - error: insert failure
//
// Pebble flow: paired with EnsureDevUser before penalty debits or top-ups in local env.
func EnsureDevWallet(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO wallets (user_id) VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING`, userID)
	return err
}

// GetUserByID loads a full user profile including streak fields and penalty settings.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: users.id primary key
//
// Returns:
//   - *models.User: populated row, or nil when not found
//   - error: database errors other than no rows
//
// Pebble flow: api-gateway GET /users/me, penalty-service rate lookup, scoring-service streak job.
func GetUserByID(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*models.User, error) {
	var u models.User
	var streakUpdated *time.Time
	err := pool.QueryRow(ctx, `
		SELECT id, email, phone, risk_profile, penalty_rate, penalty_threshold,
		       invest_threshold, consent_hours, streak_count, streak_last_updated,
		       created_at, updated_at
		FROM users WHERE id = $1`, userID,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.RiskProfile, &u.PenaltyRate, &u.PenaltyThreshold,
		&u.InvestThreshold, &u.ConsentHours, &u.StreakCount, &streakUpdated,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.StreakLastUpdated = streakUpdated
	return &u, nil
}

// WeekImpulseStats holds rolling 7-day impulse aggregates used for streak decisions.
type WeekImpulseStats struct {
	AvgScore      float64
	TxCount       int
	ImpulsePct    float64 // share of line items above user threshold
}

// GetWeekImpulseStats computes 7-day impulse behaviour for streak evaluation.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: transactions.user_id filter
//   - threshold: penalty_threshold; scores >= threshold count toward ImpulsePct
//
// Returns:
//   - *WeekImpulseStats: average score, distinct transaction count, impulse line-item %
//   - error: query failure
//
// Pebble flow: scoring-service weekly cron compares ImpulsePct against streak rules
// before calling IncrementStreak and publishing queue.TopicStreakUpdated.
func GetWeekImpulseStats(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, threshold int) (*WeekImpulseStats, error) {
	var stats WeekImpulseStats
	err := pool.QueryRow(ctx, `
		SELECT COALESCE(AVG(li.impulse_score), 0),
		       COUNT(DISTINCT t.id),
		       COALESCE(
		         100.0 * SUM(CASE WHEN li.impulse_score >= $2 THEN 1 ELSE 0 END)::float
		         / NULLIF(COUNT(li.id), 0),
		         0
		       )
		FROM transactions t
		JOIN line_items li ON li.transaction_id = t.id
		WHERE t.user_id = $1
		  AND t.logged_at >= NOW() - INTERVAL '7 days'`,
		userID, threshold,
	).Scan(&stats.AvgScore, &stats.TxCount, &stats.ImpulsePct)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// IncrementStreak bumps users.streak_count and sets streak_last_updated to now.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: user earning the streak week
//
// Returns:
//   - int: new streak_count after increment
//   - error: update failure
//
// Pebble flow: scoring-service after a qualifying low-impulse week; each +4 streak weeks
// reduces effective penalty rate via EffectivePenaltyRateForUser.
func IncrementStreak(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		UPDATE users
		SET streak_count = streak_count + 1,
		    streak_last_updated = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING streak_count`, userID,
	).Scan(&count)
	return count, err
}
