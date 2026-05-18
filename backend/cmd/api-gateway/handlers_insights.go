// Package main (handlers_insights.go) serves weekly digest and peer-benchmark insights
// aggregated from scored transactions and penalties in PostgreSQL (no RabbitMQ).
package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
)

// handleWeeklyInsights serves GET /api/v1/insights/weekly for the authenticated user's digest.
func handleWeeklyInsights(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		user, err := queries.GetUserByID(r.Context(), db, userID)
		if err != nil || user == nil {
			httputil.RespondError(w, http.StatusNotFound, "user not found", "NOT_FOUND")
			return
		}
		digest, err := queries.GetWeeklyDigest(r.Context(), db, userID, user.PenaltyThreshold)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to build digest", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, digest)
	}
}

// handleBenchmarkInsights serves GET /api/v1/insights/benchmark comparing the user to cohort peers.
func handleBenchmarkInsights(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		user, err := queries.GetUserByID(r.Context(), db, userID)
		if err != nil || user == nil {
			httputil.RespondError(w, http.StatusNotFound, "user not found", "NOT_FOUND")
			return
		}
		bench, err := queries.GetBenchmark(r.Context(), db, userID, user.PenaltyThreshold)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to build benchmark", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, bench)
	}
}
