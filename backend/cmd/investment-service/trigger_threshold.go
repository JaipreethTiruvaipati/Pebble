// Package main (trigger_threshold.go) runs the pooled-balance threshold trigger: when total
// confirmed penalties in the pool reach DefaultPoolThresholdINR, ExecutePool("threshold") runs.
package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/rs/zerolog/log"
)

// DefaultPoolThresholdINR is the minimum pooled balance before the threshold trigger fires.
const DefaultPoolThresholdINR = 500.0

// StartThresholdTrigger polls SumPooledAmount every 5 minutes and calls ExecutePool("threshold") when sum ≥ DefaultPoolThresholdINR.
func StartThresholdTrigger(db *pgxpool.Pool) {
	log.Info().Float64("threshold_inr", DefaultPoolThresholdINR).Msg("started threshold trigger goroutine")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sum, err := queries.SumPooledAmount(context.Background(), db)
		if err != nil {
			log.Error().Err(err).Msg("threshold trigger: pool sum query failed")
			continue
		}
		if sum >= DefaultPoolThresholdINR {
			log.Info().Float64("pooled_total", sum).Msg("threshold exceeded — executing pool")
			ExecutePool("threshold")
		}
	}
}
