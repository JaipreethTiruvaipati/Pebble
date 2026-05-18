package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("service", "investment-service").Logger()
	log.Info().Msg("starting investment-service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	_ = cfg // Phase 2: db pool, redis, broker client init

	// Setup Prometheus metrics endpoint
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Info().Msg("metrics server listening on :9095")
		if err := http.ListenAndServe(":9095", mux); err != nil {
			log.Error().Err(err).Msg("metrics server failed")
		}
	}()

	// Start asynchronous background triggers
	go StartThresholdTrigger()
	go StartTimeTrigger()
	go StartOpportunityTrigger()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down investment-service...")
}
