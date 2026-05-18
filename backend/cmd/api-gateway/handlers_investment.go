package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/auth"
	"github.com/jaipreeth/pebble/backend/internal/cache"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
	"github.com/jaipreeth/pebble/backend/internal/models"
)

func userIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	v := r.Context().Value(auth.UserIDKey)
	id, ok := v.(uuid.UUID)
	return id, ok
}

func handleGetPortfolio(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		summary, err := queries.GetPortfolioSummary(r.Context(), db, userID)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to load portfolio", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, summary)
	}
}

func handleListInvestments(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		triggerType := r.URL.Query().Get("trigger_type")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		items, err := queries.ListInvestments(r.Context(), db, userID, triggerType, limit)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to list investments", "INTERNAL_ERROR")
			return
		}
		if items == nil {
			items = []models.Investment{}
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"investments": items,
			"total":       len(items),
		})
	}
}

func handleGetInvestment(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid investment id", "INVALID_ID")
			return
		}
		inv, err := queries.GetInvestmentByID(r.Context(), db, userID, id)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to load investment", "INTERNAL_ERROR")
			return
		}
		if inv == nil {
			httputil.RespondError(w, http.StatusNotFound, "investment not found", "NOT_FOUND")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, inv)
	}
}

func handleGetMarketSignal(redis *cache.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var signals []models.MarketSignal
		var updatedAt string
		found := false
		if redis != nil {
			ok, err := redis.GetJSON(r.Context(), cache.KeyMarketSignals, &signals)
			if err != nil {
				httputil.RespondError(w, http.StatusInternalServerError, "failed to read market signals", "INTERNAL_ERROR")
				return
			}
			found = ok
		}
		composite := 0.0
		for _, s := range signals {
			if s.Action == "BUY" {
				composite += s.Value
			}
		}
		if len(signals) > 0 {
			composite /= float64(len(signals))
		}
		if found {
			updatedAt = "cached"
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"signals":          signals,
			"composite_score":  composite,
			"updated_at":       updatedAt,
		})
	}
}
