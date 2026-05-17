package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/db"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("starting penalty-service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbPool, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer dbPool.Close()

	rmq, err := queue.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to rabbitmq")
	}
	defer rmq.Close()

	// 1. ProcessExpiredPending Goroutine
	// Runs every 5 minutes to automatically confirm penalties whose 24h consent window has expired.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			<-ticker.C
			log.Info().Msg("running ProcessExpiredPending job...")
			// TODO: Phase 2 - Query DB for status='pending' AND expires_at < NOW()
			// Update status to 'confirmed'
			// Publish "wallet.penalties_confirmed" event via rmq.Publish
		}
	}()

	// 2. Consume Bills Scored event to calculate new penalties
	err = rmq.Consume("penalty.bills.scored", queue.TopicBillsScored, func(body []byte) error {
		log.Info().Msg("received bills.scored event, calculating penalties...")
		
		// TODO: Phase 2 implementation:
		// 1. Fetch LineItems for this bill
		// 2. Fetch User's PenaltyThreshold and PenaltyRate
		// 3. For each item: if ImpulseScore >= PenaltyThreshold -> CalculatePenalty()
		// 4. Create pending Penalty in DB (ExpiresAt = NOW + 24 hours)
		// 5. Publish "wallet.penalty_queued" event via rmq.Publish so frontend updates pending total
		
		return nil
	})

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start consumer")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down penalty-service...")
}
