package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type BillUploadedEvent struct {
	S3Key  string `json:"s3_key"`
	UserID string `json:"user_id"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("service", "scoring-service").Logger()

	log.Info().Msg("starting scoring-service...")

	// Setup Prometheus metrics endpoint
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Info().Msg("metrics server listening on :9091")
		if err := http.ListenAndServe(":9091", mux); err != nil {
			log.Error().Err(err).Msg("metrics server failed")
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Setup RabbitMQ
	rmq, err := queue.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to rabbitmq")
	}
	defer rmq.Close()

	// Consume messages
	err = rmq.Consume("scoring.bills.uploaded", "bills.uploaded", func(body []byte) error {
		var event BillUploadedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		
		log.Info().Str("s3_key", event.S3Key).Msg("received bill uploaded event")
		
		// Pipeline:
		// 1. Download image from S3 (using event.S3Key)
		// 2. Pass to Google Vision API -> get raw OCR text
		// 3. Pass raw text to Gemini ExtractAndScore
		// 4. Save structured transaction & ScoredItems to Postgres
		// 5. Cache summary in Redis
		// 6. Publish 'bills.scored' event to queue
		
		return nil
	})

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start consumer")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down scoring-service...")
}
