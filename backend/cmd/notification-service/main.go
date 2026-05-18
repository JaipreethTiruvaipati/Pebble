// Package main runs the notification-service microservice: the fan-out step for user-facing
// alerts. It consumes wallet.penalty_queued, investments.executed, and streak.updated from
// RabbitMQ and will dispatch FCM push and SES email (Phase 2). No HTTP API except Prometheus :9093.
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// main connects to RabbitMQ, registers three consumers (penalty queued, investments executed,
// streak updated), serves /metrics on :9093, and blocks until SIGINT/SIGTERM.
func main() {
	// Structured zerolog logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("service", "notification-service").Logger()

	log.Info().Msg("starting notification-service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Setup Prometheus metrics endpoint
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Info().Msg("metrics server listening on :9093")
		if err := http.ListenAndServe(":9093", mux); err != nil {
			log.Error().Err(err).Msg("metrics server failed")
		}
	}()

	rmq, err := queue.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to rabbitmq")
	}
	defer rmq.Close()

	// consumePenaltyQueued handles wallet.penalty_queued from penalty-service (pending consent alert).
	err = rmq.Consume("notification.penalty.queued", queue.TopicWalletPenaltyQueued, func(body []byte) error {
		var event queue.PenaltyQueuedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		
		log.Info().
			Str("user_id", event.UserID.String()).
			Float64("total_pending", event.TotalPending).
			Msg("dispatching push notification via FCM and email via SES")
			
		// TODO: Phase 2 - Initialize Firebase Admin SDK and send FCM push notification
		// TODO: Phase 2 - Initialize AWS SES client and send email
		
		return nil
	})

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start penalty queued consumer")
	}

	// consumeInvestmentsExecuted handles investments.executed from investment-service.
	err = rmq.Consume("notification.investments.executed", queue.TopicInvestmentsExecuted, func(body []byte) error {
		var event queue.InvestmentExecutedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		log.Info().
			Str("trigger", event.TriggerType).
			Float64("total", event.TotalAmount).
			Str("broker_ref", event.BrokerRef).
			Int("investments", len(event.InvestmentIDs)).
			Msg("dispatching investment confirmation via FCM and SES")
		// TODO: FCM — "Rs X invested across equity, gold, bonds"
		// TODO: SES — investment receipt email
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start investments executed consumer")
	}

	// consumeStreakUpdated handles streak.updated from scoring-service after a low-impulse week.
	_ = rmq.Consume("notification.streak.updated", queue.TopicStreakUpdated, func(body []byte) error {
		var event queue.StreakUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		log.Info().
			Int("streak", event.StreakCount).
			Str("user_id", event.UserID.String()).
			Msg("streak milestone — notify user of discipline streak")
		return nil
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down notification-service...")
}
