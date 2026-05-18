// Package main (trigger_time.go) runs the monthly SIP time trigger on the 1st at 9:00 AM IST,
// guaranteeing periodic investment even when the pool threshold is not met.
package main

import (
	"time"

	"github.com/rs/zerolog/log"
)

// ShouldRunMonthlySIP reports whether now (in loc) is the 1st of the month at hour 9.
func ShouldRunMonthlySIP(now time.Time, loc *time.Location) bool {
	t := now.In(loc)
	return t.Day() == 1 && t.Hour() == 9
}

// StartTimeTrigger checks hourly and calls ExecutePool("time") when ShouldRunMonthlySIP is true in Asia/Kolkata.
func StartTimeTrigger() {
	log.Info().Msg("started time trigger goroutine (1st of month, 9 AM IST)")
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		ist = time.FixedZone("IST", 5*3600+1800)
	}

	for range ticker.C {
		if ShouldRunMonthlySIP(time.Now(), ist) {
			log.Info().Msg("monthly SIP window — executing pool")
			ExecutePool("time")
		}
	}
}
