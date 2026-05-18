// Package main runs the bill-service microservice: the ingestion step after a user uploads
// a receipt. It accepts multipart POST /upload, stores images in S3, and publishes
// bills.uploaded on RabbitMQ (pebble.events exchange) for scoring-service to OCR and score.
package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// main connects to RabbitMQ, wires the S3 uploader, registers POST /upload on port 8081,
// and blocks until the HTTP server exits.
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("service", "bill-service").Logger()
	log.Info().Msg("starting bill-service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	rmq, err := queue.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to rabbitmq")
	}
	defer rmq.Close()

	uploader := NewS3Uploader(cfg.AWSS3Bucket, cfg.AWSRegion)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/upload", HandleBillUpload(uploader, rmq))

	log.Info().Msg("bill-service listening on port 8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
