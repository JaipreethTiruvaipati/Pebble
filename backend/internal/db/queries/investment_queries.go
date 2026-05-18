// Package queries encapsulates parameterized SQL access for Pebble domain tables.
// Investment queries support portfolio APIs and the investment-service batch executor
// that converts pooled penalty cash into per-user holdings.
package queries

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/models"
)

// PortfolioSummary aggregates a user's active investments by asset class for the dashboard.
// GainPct is a placeholder until live NAV feeds are wired; Allocation maps class → % of total.
type PortfolioSummary struct {
	TotalInvested float64            `json:"total_invested"`
	EquityValue   float64            `json:"equity_value"`
	GoldValue     float64            `json:"gold_value"`
	BondValue     float64            `json:"bond_value"`
	GainPct       float64            `json:"gain_pct"`
	Allocation    map[string]float64 `json:"allocation_pct"`
}

// ListInvestments returns a user's investment rows, optionally filtered by trigger type.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: owner of the investments table rows
//   - triggerType: when non-empty, filters investments.trigger_type (e.g. "threshold", "market_signal")
//   - limit: max rows; values <= 0 default to 20
//
// Returns:
//   - []models.Investment: newest first (created_at DESC)
//   - error: query or scan failure
//
// Pebble flow: api-gateway GET /portfolio/investments reads this after JWT auth; rows map
// to confirmed broker purchases created by MarkPoolInvested or manual top-up flows.
func ListInvestments(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, triggerType string, limit int) ([]models.Investment, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, user_id, asset_class, amount, units, nav_at_purchase, status, created_at,
		       COALESCE(trigger_type, ''), COALESCE(broker_ref, '')
		FROM investments
		WHERE user_id = $1`
	args := []interface{}{userID}
	if triggerType != "" {
		query += ` AND trigger_type = $2`
		args = append(args, triggerType)
	}
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d`, limit)

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Investment
	for rows.Next() {
		var inv models.Investment
		var trigger, brokerRef string
		if err := rows.Scan(
			&inv.ID, &inv.UserID, &inv.AssetClass, &inv.Amount, &inv.Units,
			&inv.NAVAtPurchase, &inv.Status, &inv.CreatedAt, &trigger, &brokerRef,
		); err != nil {
			return nil, err
		}
		if trigger != "" {
			inv.TriggerType = trigger
		}
		if brokerRef != "" {
			inv.BrokerRef = brokerRef
		}
		out = append(out, inv)
	}
	return out, rows.Err()
}

// GetPortfolioSummary computes per-asset-class totals and allocation percentages for one user.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: investments.user_id filter
//
// Returns:
//   - *PortfolioSummary: aggregated values and allocation_pct map (empty when no holdings)
//   - error: query or scan failure
//
// How it works: SUM(amount) GROUP BY asset_class for status='active', buckets into equity/gold/bonds,
// then derives allocation percentages. Used by api-gateway portfolio summary endpoint.
func GetPortfolioSummary(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*PortfolioSummary, error) {
	rows, err := pool.Query(ctx, `
		SELECT asset_class, COALESCE(SUM(amount), 0)
		FROM investments
		WHERE user_id = $1 AND status = 'active'
		GROUP BY asset_class`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := &PortfolioSummary{Allocation: make(map[string]float64)}
	for rows.Next() {
		var class string
		var amt float64
		if err := rows.Scan(&class, &amt); err != nil {
			return nil, err
		}
		summary.TotalInvested += amt
		switch class {
		case "equity":
			summary.EquityValue = amt
		case "gold":
			summary.GoldValue = amt
		case "bonds", "bond":
			summary.BondValue += amt
		default:
			// mutual_funds and others count toward total only
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if summary.TotalInvested > 0 {
		if summary.EquityValue > 0 {
			summary.Allocation["equity"] = (summary.EquityValue / summary.TotalInvested) * 100
		}
		if summary.GoldValue > 0 {
			summary.Allocation["gold"] = (summary.GoldValue / summary.TotalInvested) * 100
		}
		if summary.BondValue > 0 {
			summary.Allocation["bonds"] = (summary.BondValue / summary.TotalInvested) * 100
		}
		// Placeholder until live NAV feed — demonstrates API shape for frontend.
		summary.GainPct = 4.2
	}
	return summary, nil
}

// SumPooledAmount returns the total INR sitting in pool_contributions awaiting batch investment.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//
// Returns:
//   - float64: COALESCE(SUM(amount), 0) where status='pooled'
//   - error: query failure
//
// Pebble flow: investment-service checks this before ExecutePool; when sum >= user's
// invest_threshold (aggregated), MarkPoolInvested runs broker allocation.
func SumPooledAmount(ctx context.Context, pool *pgxpool.Pool) (float64, error) {
	var sum float64
	err := pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0) FROM pool_contributions WHERE status = 'pooled'`).Scan(&sum)
	return sum, err
}

// MarkPoolInvested atomically allocates pooled penalty cash into per-user investments.
//
// Parameters:
//   - ctx: transaction context
//   - pool: shared PostgreSQL pool (begins a local tx)
//   - triggerType: stored on each new investments row (e.g. "threshold")
//   - brokerRef: external order reference from Smallcase/broker
//   - splits: map of asset_class → total INR to deploy across the whole pool this batch
//
// Returns:
//   - []uuid.UUID: IDs of inserted investments rows (one per user × asset class slice)
//   - nil, nil when no pooled contributions exist
//   - error: tx, insert, or update failure (rolled back)
//
// How it works: loads all pool_contributions with status='pooled', pro-rates each user's
// share of splits by contribution amount, inserts investments with placeholder NAV/units,
// then marks contributions invested. After commit, investment-service publishes
// queue.TopicInvestmentsExecuted with the returned IDs.
func MarkPoolInvested(ctx context.Context, pool *pgxpool.Pool, triggerType, brokerRef string, splits map[string]float64) ([]uuid.UUID, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT id, user_id, amount FROM pool_contributions WHERE status = 'pooled'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type contrib struct {
		id     uuid.UUID
		userID uuid.UUID
		amount float64
	}
	var contribs []contrib
	var totalPool float64
	for rows.Next() {
		var c contrib
		if err := rows.Scan(&c.id, &c.userID, &c.amount); err != nil {
			return nil, err
		}
		contribs = append(contribs, c)
		totalPool += c.amount
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if totalPool <= 0 {
		return nil, nil
	}

	var investmentIDs []uuid.UUID
	for _, c := range contribs {
		share := c.amount / totalPool
		for assetClass, splitTotal := range splits {
			if splitTotal <= 0 {
				continue
			}
			amt := splitTotal * share
			if amt < 0.01 {
				continue
			}
			nav := 100.0
			if assetClass == "gold" {
				nav = 6500.0
			}
			units := amt / nav
			var invID uuid.UUID
			err := tx.QueryRow(ctx, `
				INSERT INTO investments (user_id, asset_class, amount, units, nav_at_purchase, status, trigger_type, broker_ref)
				VALUES ($1, $2, $3, $4, $5, 'active', $6, $7)
				RETURNING id`,
				c.userID, assetClass, amt, units, nav, triggerType, brokerRef,
			).Scan(&invID)
			if err != nil {
				return nil, err
			}
			investmentIDs = append(investmentIDs, invID)
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE pool_contributions SET status = 'invested', invested_at = NOW() WHERE status = 'pooled'`)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return investmentIDs, nil
}

// GetInvestmentByID loads one investment row scoped to the authenticated user.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: row owner (authorization guard)
//   - investmentID: primary key
//
// Returns:
//   - *models.Investment: populated row, or nil when not found
//   - error: database errors other than no rows
func GetInvestmentByID(ctx context.Context, pool *pgxpool.Pool, userID, investmentID uuid.UUID) (*models.Investment, error) {
	var inv models.Investment
	var trigger, brokerRef string
	err := pool.QueryRow(ctx, `
		SELECT id, user_id, asset_class, amount, units, nav_at_purchase, status, created_at,
		       COALESCE(trigger_type, ''), COALESCE(broker_ref, '')
		FROM investments WHERE id = $1 AND user_id = $2`,
		investmentID, userID,
	).Scan(
		&inv.ID, &inv.UserID, &inv.AssetClass, &inv.Amount, &inv.Units,
		&inv.NAVAtPurchase, &inv.Status, &inv.CreatedAt, &trigger, &brokerRef,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	inv.TriggerType = trigger
	inv.BrokerRef = brokerRef
	return &inv, nil
}
