// Package main runs the notification-service microservice: the fan-out step for user-facing
// alerts. It consumes wallet.penalty_queued, investments.executed, and streak.updated from
// RabbitMQ and dispatches FCM push and SES email via the Dispatcher.
// No HTTP API except Prometheus :9093.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/jaipreeth/pebble/backend/pkg/notify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// main connects to RabbitMQ, initializes FCM and SES clients, registers three consumers
// (penalty queued, investments executed, streak updated), serves /metrics on :9093,
// and blocks until SIGINT/SIGTERM.
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

	// Initialize FCM client (optional — skips push if credentials unavailable)
	ctx := context.Background()
	var fcmClient *notify.FCMClient
	if cfg.FirebaseCredPath != "" {
		fcmClient, err = notify.NewFCMClient(ctx, cfg.FirebaseCredPath)
		if err != nil {
			log.Warn().Err(err).Msg("FCM client unavailable — push notifications disabled")
		} else {
			log.Info().Msg("FCM client initialized")
		}
	}

	// Initialize SES client (optional — skips email if not configured)
	var sesClient *notify.SESClient
	if cfg.SESFromEmail != "" {
		sesClient, err = notify.NewSESClient(ctx, cfg.AWSRegion, cfg.SESFromEmail)
		if err != nil {
			log.Warn().Err(err).Msg("SES client unavailable — email notifications disabled")
		} else {
			log.Info().Msg("SES client initialized")
		}
	}

	dispatcher := NewDispatcher(fcmClient, sesClient)

	// consumePenaltyQueued handles wallet.penalty_queued from penalty-service (pending consent alert).
	err = rmq.Consume("notification.penalty.queued", queue.TopicWalletPenaltyQueued, func(body []byte) error {
		var event queue.PenaltyQueuedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Info().
			Str("user_id", event.UserID.String()).
			Float64("total_pending", event.TotalPending).
			Msg("dispatching penalty notification via FCM and SES")

		dispatcher.NotifyPenaltyQueued(ctx, event.UserID, event.TotalPending)
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

		dispatcher.NotifyInvestmentExecuted(ctx, event.TriggerType, event.TotalAmount, event.BrokerRef)
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

		dispatcher.NotifyStreakMilestone(ctx, event.UserID, event.StreakCount)
		return nil
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down notification-service...")
}
