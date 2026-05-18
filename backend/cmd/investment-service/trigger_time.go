package main

import (
	"time"

	"github.com/rs/zerolog/log"
)

// StartTimeTrigger ensures funds don't sit idle for too long.
// It executes the pool on the 1st of every month at 9:00 AM IST.
func StartTimeTrigger() {
	log.Info().Msg("started time trigger goroutine (1st of month)")
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Load India Standard Time timezone
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		// Fallback to manual offset
		ist = time.FixedZone("IST", 5*3600+1800)
	}

	for range ticker.C {
		now := time.Now().In(ist)
		// Check if it is the 1st day of the month and the hour is 9 AM
		if now.Day() == 1 && now.Hour() == 9 {
			ExecutePool("monthly_time_trigger")
		}
	}
}
