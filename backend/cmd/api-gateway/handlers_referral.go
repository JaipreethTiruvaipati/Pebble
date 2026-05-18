// Package main (handlers_referral.go) manages referral codes: users share codes for a 2%
// penalty-rate discount (applied by penalty-service via HasReferrerDiscount).
package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
)

// handleGetReferralMe serves GET /api/v1/referrals/me (personal code and redemption stats).
func handleGetReferralMe(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		stats, err := queries.GetReferralStats(r.Context(), db, userID)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to load referral stats", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, stats)
	}
}

// handleRedeemReferral serves POST /api/v1/referrals/redeem with JSON body {"code":"..."}.
func handleRedeemReferral(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		var req struct {
			Code string `json:"code"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid JSON body", "INVALID_REQUEST")
			return
		}
		code, err := httputil.ValidateReferralCode(req.Code)
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, err.Error(), "INVALID_REQUEST")
			return
		}
		if err := queries.RedeemReferralCode(r.Context(), db, userID, code); err != nil {
			switch {
			case errors.Is(err, queries.ErrReferralInvalidCode):
				httputil.RespondError(w, http.StatusBadRequest, "invalid referral code", "INVALID_CODE")
			case errors.Is(err, queries.ErrReferralSelfRedeem):
				httputil.RespondError(w, http.StatusBadRequest, "cannot redeem your own code", "SELF_REDEEM")
			case errors.Is(err, queries.ErrReferralAlreadyRedeemed):
				httputil.RespondError(w, http.StatusConflict, "referral already redeemed", "ALREADY_REDEEMED")
			default:
				httputil.RespondError(w, http.StatusInternalServerError, "failed to redeem code", "INTERNAL_ERROR")
			}
			return
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]string{
			"message": "referral code redeemed successfully",
		})
	}
}
