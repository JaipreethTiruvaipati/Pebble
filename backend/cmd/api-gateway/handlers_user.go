// Package main (handlers_user.go) serves the authenticated user profile including streak
// metadata (updated when scoring-service publishes streak.updated after a low-impulse week).
package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
)

// handleGetMe serves GET /api/v1/me and GET /api/v1/users/me with effective penalty rate
// (base rate minus streak and referral discounts, matching penalty-service logic).
func handleGetMe(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}

		user, err := queries.GetUserByID(r.Context(), db, userID)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to load profile", "INTERNAL_ERROR")
			return
		}
		if user == nil {
			_ = queries.EnsureDevUser(r.Context(), db, userID, "")
			_ = queries.EnsureDevWallet(r.Context(), db, userID)
			user, _ = queries.GetUserByID(r.Context(), db, userID)
		}
		if user == nil {
			httputil.RespondError(w, http.StatusNotFound, "user not found", "NOT_FOUND")
			return
		}

		hasReferral, _ := queries.HasReferrerDiscount(r.Context(), db, userID)
		effectiveRate := queries.EffectivePenaltyRateForUser(user.PenaltyRate, user.StreakCount, hasReferral)
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"id":                  user.ID,
			"email":               user.Email,
			"phone":               user.Phone,
			"risk_profile":        user.RiskProfile,
			"penalty_rate":        user.PenaltyRate,
			"effective_penalty_rate": effectiveRate,
			"penalty_threshold":   user.PenaltyThreshold,
			"invest_threshold":    user.InvestThreshold,
			"consent_hours":       user.ConsentHours,
			"streak_count":        user.StreakCount,
			"streak_last_updated": user.StreakLastUpdated,
			"streak_discount_pct": float64(user.StreakCount/4) * 1.0,
		})
	}
}
