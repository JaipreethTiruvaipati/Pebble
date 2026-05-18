package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
)

func handleGetWalletBalance(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		_ = queries.EnsureDevWallet(r.Context(), db, userID)
		snap, err := queries.GetWalletSnapshot(r.Context(), db, userID)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to load wallet", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, snap)
	}
}

func handleGetWalletLedger(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		entries, err := queries.ListWalletLedger(r.Context(), db, userID, 50)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to load ledger", "INTERNAL_ERROR")
			return
		}
		if entries == nil {
			entries = []queries.WalletLedgerEntry{}
		}
		httputil.RespondJSON(w, http.StatusOK, entries)
	}
}

func handleWalletTopup(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		var req struct {
			Amount float64 `json:"amount"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Amount <= 0 {
			req.Amount = 5000
		}
		_ = queries.EnsureDevWallet(r.Context(), db, userID)
		if err := queries.CreditWalletTopup(r.Context(), db, userID, req.Amount); err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "topup failed", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"message": "topup successful",
			"amount":  req.Amount,
		})
	}
}

func handleListPenalties(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		status := r.URL.Query().Get("status")
		items, err := queries.ListPenalties(r.Context(), db, userID, status)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to list penalties", "INTERNAL_ERROR")
			return
		}
		if items == nil {
			items = []queries.PenaltyRow{}
		}
		httputil.RespondJSON(w, http.StatusOK, items)
	}
}

func handleContestPenalty(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		pid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid id", "INVALID_ID")
			return
		}
		if err := queries.ContestPenalty(r.Context(), db, userID, pid); err != nil {
			if err == pgx.ErrNoRows {
				httputil.RespondError(w, http.StatusNotFound, "penalty not found", "NOT_FOUND")
				return
			}
			httputil.RespondError(w, http.StatusInternalServerError, "failed to contest", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"id": pid.String(), "status": "contested"})
	}
}

func handleConfirmPenaltyEarly(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		pid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid id", "INVALID_ID")
			return
		}
		if err := queries.ConfirmPenaltyEarly(r.Context(), db, userID, pid); err != nil {
			if err == pgx.ErrNoRows {
				httputil.RespondError(w, http.StatusNotFound, "penalty not found", "NOT_FOUND")
				return
			}
			httputil.RespondError(w, http.StatusInternalServerError, "failed to confirm", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"id": pid.String(), "status": "confirmed"})
	}
}

func handleGetPendingPenalty(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		b, err := queries.GetTopPendingPenalty(r.Context(), db, userID)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to load penalty", "INTERNAL_ERROR")
			return
		}
		if b == nil {
			httputil.RespondJSON(w, http.StatusOK, nil)
			return
		}
		httputil.RespondJSON(w, http.StatusOK, b)
	}
}
