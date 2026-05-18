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

// StartThresholdTrigger monitors pooled penalty cash and executes when the sum crosses the threshold.
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
