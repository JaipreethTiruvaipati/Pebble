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

	// Consume Penalty Queued events
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

	// Consume investment executed events (Week 17)
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down notification-service...")
}
