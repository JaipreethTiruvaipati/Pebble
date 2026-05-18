package main

import (
	"context"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/cache"
	"github.com/jaipreeth/pebble/backend/internal/models"
	"github.com/rs/zerolog/log"
)

// HasStrongBuySignal returns true when any cached signal is a high-confidence BUY.
func HasStrongBuySignal(signals []models.MarketSignal) bool {
	for _, s := range signals {
		if s.Action == "BUY" && s.Value >= 70 {
			return true
		}
	}
	return false
}

// StartOpportunityTrigger deploys the pool when Redis market signals indicate a strong opportunity.
func StartOpportunityTrigger(redis *cache.Client) {
	log.Info().Msg("started opportunity trigger goroutine (hourly signal check)")
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		var signals []models.MarketSignal
		if redis != nil {
			found, err := redis.GetJSON(context.Background(), cache.KeyMarketSignals, &signals)
			if err != nil {
				log.Error().Err(err).Msg("opportunity trigger: redis read failed")
				continue
			}
			if !found || !HasStrongBuySignal(signals) {
				continue
			}
		} else {
			continue
		}
		log.Info().Msg("strong market opportunity — executing pool")
		ExecutePool("opportunity")
	}
}
