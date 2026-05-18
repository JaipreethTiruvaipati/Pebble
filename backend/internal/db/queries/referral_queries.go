// Package queries encapsulates parameterized SQL access for Pebble domain tables.
// Referral queries manage invite codes, redemptions, and referrer penalty discounts.
package queries

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/models"
)

// ReferrerDiscountPct is the penalty-rate reduction (2%) applied when a user has ≥1 successful referral.
const ReferrerDiscountPct = 0.02 // -2% penalty rate per active referral programme

// ErrReferralSelfRedeem is returned when a user attempts to redeem their own code.
var ErrReferralSelfRedeem = errors.New("cannot redeem own referral code")

// ErrReferralAlreadyRedeemed is returned when referred_user_id already has a redemption row.
var ErrReferralAlreadyRedeemed = errors.New("user already redeemed a referral code")

// ErrReferralInvalidCode is returned when the code does not exist in referral_codes.
var ErrReferralInvalidCode = errors.New("invalid referral code")

// EnsureReferralCode creates a referral_codes row for the user if one does not exist.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: code owner
//
// Returns:
//   - *models.ReferralCode: existing or newly inserted code (format PEBBLE-{uuid-prefix})
//   - error: lookup or insert failure
//
// Pebble flow: api-gateway GET /referrals/me lazily provisions a shareable code.
func EnsureReferralCode(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*models.ReferralCode, error) {
	existing, err := GetReferralCodeByUser(ctx, pool, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	code := fmt.Sprintf("PEBBLE-%s", strings.ToUpper(userID.String()[:8]))
	var rc models.ReferralCode
	err = pool.QueryRow(ctx, `
		INSERT INTO referral_codes (user_id, code)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET code = referral_codes.code
		RETURNING id, user_id, code, created_at`,
		userID, code,
	).Scan(&rc.ID, &rc.UserID, &rc.Code, &rc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rc, nil
}

// GetReferralCodeByUser returns the referral_codes row for a user, if any.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: referral_codes.user_id
//
// Returns:
//   - *models.ReferralCode: row or nil when not found
//   - error: database errors other than no rows
func GetReferralCodeByUser(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*models.ReferralCode, error) {
	var rc models.ReferralCode
	err := pool.QueryRow(ctx, `
		SELECT id, user_id, code, created_at FROM referral_codes WHERE user_id = $1`, userID,
	).Scan(&rc.ID, &rc.UserID, &rc.Code, &rc.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rc, nil
}

// CountReferralRedemptions counts how many users redeemed the referrer's code.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - referrerUserID: owner of the referral_codes row (not the referred user)
//
// Returns:
//   - int: COUNT of referral_redemptions joined to this user's code
//   - error: query failure
func CountReferralRedemptions(ctx context.Context, pool *pgxpool.Pool, referrerUserID uuid.UUID) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM referral_redemptions rr
		JOIN referral_codes rc ON rc.id = rr.referral_code_id
		WHERE rc.user_id = $1`, referrerUserID,
	).Scan(&count)
	return count, err
}

// HasReferrerDiscount reports whether the user qualifies for the referrer penalty discount.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: potential referrer (code owner)
//
// Returns:
//   - bool: true when CountReferralRedemptions > 0
//   - error: count query failure
//
// Pebble flow: penalty-service passes this into EffectivePenaltyRateForUser alongside streak.
func HasReferrerDiscount(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (bool, error) {
	n, err := CountReferralRedemptions(ctx, pool, userID)
	return n > 0, err
}

// RedeemReferralCode records a new user's redemption of someone else's invite code.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - referredUserID: the new user applying the code (one redemption per user)
//   - code: invite string (trimmed and uppercased)
//
// Returns:
//   - nil on success
//   - ErrReferralInvalidCode, ErrReferralSelfRedeem, ErrReferralAlreadyRedeemed, or DB error
//
// Pebble flow: api-gateway POST /referrals/redeem after signup; unique constraint on
// referred_user_id enforces single redemption per account.
func RedeemReferralCode(ctx context.Context, pool *pgxpool.Pool, referredUserID uuid.UUID, code string) error {
	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return ErrReferralInvalidCode
	}

	var codeID uuid.UUID
	var ownerID uuid.UUID
	err := pool.QueryRow(ctx, `
		SELECT id, user_id FROM referral_codes WHERE UPPER(code) = $1`, code,
	).Scan(&codeID, &ownerID)
	if err == pgx.ErrNoRows {
		return ErrReferralInvalidCode
	}
	if err != nil {
		return err
	}
	if ownerID == referredUserID {
		return ErrReferralSelfRedeem
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO referral_redemptions (referral_code_id, referred_user_id)
		VALUES ($1, $2)`, codeID, referredUserID)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return ErrReferralAlreadyRedeemed
		}
		return err
	}
	return nil
}

// GetReferralStats assembles code, redemption count, and active discount for the referrals API.
//
// Parameters:
//   - ctx: query context
//   - pool: shared PostgreSQL pool
//   - userID: referrer viewing their stats
//
// Returns:
//   - *models.ReferralStats: code string, redemption_count, discount_pct (2×100 when count > 0)
//   - error: ensure or count failure
//
// Pebble flow: api-gateway GET /referrals/me; DiscountPct mirrors ReferrerDiscountPct used
// in penalty rate calculation.
func GetReferralStats(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) (*models.ReferralStats, error) {
	rc, err := EnsureReferralCode(ctx, pool, userID)
	if err != nil {
		return nil, err
	}
	count, err := CountReferralRedemptions(ctx, pool, userID)
	if err != nil {
		return nil, err
	}
	discount := 0.0
	if count > 0 {
		discount = ReferrerDiscountPct * 100
	}
	return &models.ReferralStats{
		Code:            rc.Code,
		RedemptionCount: count,
		DiscountPct:     discount,
	}, nil
}
