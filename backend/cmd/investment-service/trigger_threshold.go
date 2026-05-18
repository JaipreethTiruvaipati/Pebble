package main

import (
	"time"

	"github.com/rs/zerolog/log"
)

// StartThresholdTrigger monitors the database pool. If the total uninvested
// cash crosses a threshold (e.g. 100,000 INR), it triggers a micro-batch.
func StartThresholdTrigger() {
	log.Info().Msg("started threshold trigger goroutine")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// TODO: Query DB
		// SELECT sum(amount) FROM pool_contributions WHERE status = 'pooled'
		// If sum >= 100_000:
		// ExecutePool("threshold_exceeded")
	}
}
