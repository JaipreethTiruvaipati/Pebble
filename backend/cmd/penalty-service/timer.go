// Package main (timer.go) provides consent-window timer helpers for the penalty-service.
// The primary expiry mechanism is the 5-minute sweep in main.go (runExpirySweep);
// this file adds precision helpers for consent window calculations.
package main

import (
	"time"
)

// DefaultConsentWindowHours is the default time users have to review and cancel a penalty
// before it auto-confirms and funds move to the investment pool.
const DefaultConsentWindowHours = 24

// IsConsentExpired checks whether a penalty created at createdAt with the given consent
// window (in hours) has expired relative to now.
func IsConsentExpired(createdAt time.Time, consentHours int) bool {
	if consentHours <= 0 {
		consentHours = DefaultConsentWindowHours
	}
	expiresAt := createdAt.Add(time.Duration(consentHours) * time.Hour)
	return time.Now().After(expiresAt)
}

// TimeUntilExpiry returns the remaining duration before a penalty's consent window expires.
// Returns 0 if already expired.
func TimeUntilExpiry(createdAt time.Time, consentHours int) time.Duration {
	if consentHours <= 0 {
		consentHours = DefaultConsentWindowHours
	}
	expiresAt := createdAt.Add(time.Duration(consentHours) * time.Hour)
	remaining := time.Until(expiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
