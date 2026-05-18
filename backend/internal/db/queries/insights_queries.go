// Package queries encapsulates parameterized SQL access for Pebble domain tables.
// Insights queries power spending digests and cohort benchmarks on the dashboard.
package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WeeklyDigest is the rolling 7-day spending summary returned by GET /insights/weekly.
type WeeklyDigest struct {
	WeekStart       time.Time          `json:"week_start"`
	WeekEnd         time.Time          `json:"week_end"`
	TotalSpend      float64            `json:"total_spend"`
	ImpulsePct      float64            `json:"impulse_pct"`
	AvgImpulseScore float64            `json:"avg_impulse_score"`
	TopCategories   []CategorySpend    `json:"top_categories"`
	TrendVsLastWeek float64            `json:"trend_vs_last_week_pct"` // negative = spent less
}

// CategorySpend is spend aggregated by line_items.category for a date window.
type CategorySpend struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Pct      float64 `json:"pct"`
}

// BenchmarkResult compares one user's impulse behaviour to their risk_profile cohort.
type BenchmarkResult struct {
	UserImpulsePct    float64 `json:"user_impulse_pct"`
	CohortImpulsePct  float64 `json:"cohort_impulse_pct"`
	SavedVsCohortPct  float64 `json:"saved_vs_cohort_pct"`
	CohortLabel       string  `json:"cohort_label"`
	SampleSize        int     `json:"sample_size"`
}

// GetWeeklyDigest builds the last-7-days spending digest for a user.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: transactions.user_id filter
//   - threshold: user's penalty_threshold; line items with impulse_score >= threshold count as impulse
//
// Returns:
//   - *WeeklyDigest: totals, impulse %, top 5 categories, week-over-week trend
//   - error: query failure
//
// Pebble flow: api-gateway insights handler; joins transactions + line_items written by
// bill-service after LLM scoring (bills.scored event path).
func GetWeeklyDigest(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, threshold int) (*WeeklyDigest, error) {
	now := time.Now()
	weekStart := now.Add(-7 * 24 * time.Hour)

	var totalSpend, avgScore, impulsePct float64
	err := pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(t.total_amount), 0),
		       COALESCE(AVG(li.impulse_score), 0),
		       COALESCE(
		         100.0 * SUM(CASE WHEN li.impulse_score >= $2 THEN 1 ELSE 0 END)::float
		         / NULLIF(COUNT(li.id), 0), 0)
		FROM transactions t
		LEFT JOIN line_items li ON li.transaction_id = t.id
		WHERE t.user_id = $1 AND t.logged_at >= $3`,
		userID, threshold, weekStart,
	).Scan(&totalSpend, &avgScore, &impulsePct)
	if err != nil {
		return nil, err
	}

	var prevSpend float64
	_ = pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE user_id = $1
		  AND logged_at >= $2 - INTERVAL '7 days'
		  AND logged_at < $2`,
		userID, weekStart,
	).Scan(&prevSpend)

	trend := 0.0
	if prevSpend > 0 {
		trend = ((totalSpend - prevSpend) / prevSpend) * 100
	}

	rows, err := pool.Query(ctx, `
		SELECT COALESCE(li.category, 'other'), COALESCE(SUM(li.amount), 0)
		FROM transactions t
		JOIN line_items li ON li.transaction_id = t.id
		WHERE t.user_id = $1 AND t.logged_at >= $2
		GROUP BY li.category
		ORDER BY SUM(li.amount) DESC
		LIMIT 5`, userID, weekStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategorySpend
	for rows.Next() {
		var c CategorySpend
		if err := rows.Scan(&c.Category, &c.Amount); err != nil {
			return nil, err
		}
		if totalSpend > 0 {
			c.Pct = (c.Amount / totalSpend) * 100
		}
		categories = append(categories, c)
	}

	return &WeeklyDigest{
		WeekStart:       weekStart,
		WeekEnd:         now,
		TotalSpend:      totalSpend,
		ImpulsePct:      impulsePct,
		AvgImpulseScore: avgScore,
		TopCategories:   categories,
		TrendVsLastWeek: trend,
	}, rows.Err()
}

// GetBenchmark compares a user's 30-day impulse % to peers with the same risk_profile.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: subject user; also used to load risk_profile from users
//   - threshold: impulse line-item cutoff (penalty_threshold)
//
// Returns:
//   - *BenchmarkResult: user vs cohort impulse %, saved delta, anonymised label, sample size
//   - error: query failure
//
// Pebble flow: api-gateway GET /insights/benchmark; cohort is all users sharing
// conservative/moderate/aggressive risk_profile with scored transactions in the window.
func GetBenchmark(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, threshold int) (*BenchmarkResult, error) {
	var riskProfile string
	if err := pool.QueryRow(ctx, `SELECT risk_profile FROM users WHERE id = $1`, userID).Scan(&riskProfile); err != nil {
		return nil, err
	}

	var userPct float64
	_ = pool.QueryRow(ctx, `
		SELECT COALESCE(
		  100.0 * SUM(CASE WHEN li.impulse_score >= $2 THEN 1 ELSE 0 END)::float
		  / NULLIF(COUNT(li.id), 0), 0)
		FROM transactions t
		JOIN line_items li ON li.transaction_id = t.id
		WHERE t.user_id = $1 AND t.logged_at >= NOW() - INTERVAL '30 days'`,
		userID, threshold,
	).Scan(&userPct)

	var cohortPct float64
	var sample int
	err := pool.QueryRow(ctx, `
		SELECT COALESCE(
		  100.0 * SUM(CASE WHEN li.impulse_score >= $2 THEN 1 ELSE 0 END)::float
		  / NULLIF(COUNT(li.id), 0), 0),
		  COUNT(DISTINCT u.id)
		FROM users u
		JOIN transactions t ON t.user_id = u.id
		JOIN line_items li ON li.transaction_id = t.id
		WHERE u.risk_profile = $1
		  AND t.logged_at >= NOW() - INTERVAL '30 days'`,
		riskProfile, threshold,
	).Scan(&cohortPct, &sample)
	if err != nil {
		return nil, err
	}

	saved := cohortPct - userPct
	if saved < 0 {
		saved = 0
	}

	return &BenchmarkResult{
		UserImpulsePct:   userPct,
		CohortImpulsePct: cohortPct,
		SavedVsCohortPct: saved,
		CohortLabel:      riskProfile + " earners (anonymised)",
		SampleSize:       sample,
	}, nil
}
