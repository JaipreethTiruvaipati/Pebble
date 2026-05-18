// Package main runs the penalty-service microservice: consumes bills.scored, computes per-line
// penalties from impulse scores, queues wallet.penalty_queued for notification-service, auto-confirms
// expired consent windows, and publishes wallet.penalties_confirmed when funds move to the
// investment pool for investment-service triggers.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/db"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// main connects to PostgreSQL and RabbitMQ, starts a 5-minute expiry sweep goroutine,
// consumes bills.scored on queue penalty.bills.scored, and serves Prometheus on :9092.
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("service", "penalty-service").Logger()
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

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":9092", mux)
	}()

	// runExpirySweepLoop auto-confirms expired penalties every 5 minutes and moves funds to the pool.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			runExpirySweep(context.Background(), dbPool, rmq)
		}
	}()

	// consumeBillsScored handles bills.scored events published by scoring-service.
	err = rmq.Consume("penalty.bills.scored", queue.TopicBillsScored, func(body []byte) error {
		var event queue.BillsScoredEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		return processBillsScored(context.Background(), dbPool, rmq, event)
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start bills.scored consumer")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down penalty-service...")
}

// processBillsScored loads line items for the transaction, applies EffectivePenaltyRate and
// CalculatePenalty per item, inserts pending penalties, and publishes wallet.penalty_queued.
func processBillsScored(ctx context.Context, dbPool *pgxpool.Pool, rmq *queue.RabbitMQ, event queue.BillsScoredEvent) error {
	baseRate, threshold, consentHours, err := queries.GetUserPenaltySettings(ctx, dbPool, event.UserID)
	if err != nil {
		return err
	}
	user, _ := queries.GetUserByID(ctx, dbPool, event.UserID)
	streak := 0
	if user != nil {
		streak = user.StreakCount
	}
	hasReferral, _ := queries.HasReferrerDiscount(ctx, dbPool, event.UserID)
	rate := EffectivePenaltyRate(baseRate, streak, hasReferral)
	items, err := queries.ListLineItemsForTransaction(ctx, dbPool, event.TransactionID)
	if err != nil {
		return err
	}

	var totalPending float64
	for _, item := range items {
		amt := CalculatePenalty(item.Amount, item.ImpulseScore, threshold, rate)
		if amt <= 0 {
			continue
		}
		if _, err := queries.CreatePendingPenalty(ctx, dbPool, event.UserID, item.ID, amt, consentHours); err != nil {
			return err
		}
		totalPending += amt
	}
	if totalPending > 0 {
		return rmq.Publish(ctx, queue.TopicWalletPenaltyQueued, queue.PenaltyQueuedEvent{
			UserID:        event.UserID,
			TransactionID: event.TransactionID,
			TotalPending:  totalPending,
		})
	}
	return nil
}

// runExpirySweep confirms penalties past the user consent window and calls moveConfirmedToPool.
func runExpirySweep(ctx context.Context, dbPool *pgxpool.Pool, rmq *queue.RabbitMQ) {
	n, err := queries.ConfirmExpiredPenalties(ctx, dbPool)
	if err != nil {
		log.Error().Err(err).Msg("expiry sweep failed")
		return
	}
	if n == 0 {
		return
	}
	log.Info().Int64("confirmed", n).Msg("auto-confirmed expired penalties")
	moveConfirmedToPool(ctx, dbPool, rmq)
}

// moveConfirmedToPool records pool_contributions for confirmed penalties and publishes
// wallet.penalties_confirmed per user for investment-service pooling and notifications.
func moveConfirmedToPool(ctx context.Context, dbPool *pgxpool.Pool, rmq *queue.RabbitMQ) {
	rows, err := dbPool.Query(ctx, `
		SELECT p.id, p.user_id, p.amount
		FROM penalties p
		WHERE p.status = 'confirmed'
		  AND NOT EXISTS (
		    SELECT 1 FROM pool_contributions pc WHERE pc.penalty_id = p.id
		  )`)
	if err != nil {
		log.Error().Err(err).Msg("pool move query failed")
		return
	}
	defer rows.Close()

	type batch struct {
		ids   []uuid.UUID
		total float64
	}
	byUser := make(map[uuid.UUID]*batch)
	for rows.Next() {
		var pid, uid uuid.UUID
		var amt float64
		if err := rows.Scan(&pid, &uid, &amt); err != nil {
			log.Error().Err(err).Msg("scan penalty row")
			return
		}
		if err := queries.AddPoolContribution(ctx, dbPool, uid, &pid, amt); err != nil {
			log.Error().Err(err).Msg("add pool contribution")
			return
		}
		if byUser[uid] == nil {
			byUser[uid] = &batch{}
		}
		byUser[uid].ids = append(byUser[uid].ids, pid)
		byUser[uid].total += amt
	}
	for uid, b := range byUser {
		_ = rmq.Publish(ctx, queue.TopicWalletPenaltyConfirmed, queue.PenaltiesConfirmedEvent{
			UserID:      uid,
			PenaltyIDs:  b.ids,
			TotalAmount: b.total,
		})
	}
}
