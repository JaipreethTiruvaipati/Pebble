// Package main runs the market-poller microservice: periodically fetches NSE, MCX, CCIL, and
// AMFI reference data, computes opportunity signals, and caches them in Redis for
// investment-service triggers and api-gateway GET /market/signal. No RabbitMQ; side-effect is Redis only.
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/cache"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/pkg/market"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// main connects to Redis, polls markets immediately and every 15 minutes, serves Prometheus on :9094,
// and shuts down on SIGINT/SIGTERM.
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("service", "market-poller").Logger()
	log.Info().Msg("starting market-poller service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	redisClient, err := cache.Connect(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("redis connection failed")
	}
	defer redisClient.Close()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Info().Msg("metrics server listening on :9094")
		if err := http.ListenAndServe(":9094", mux); err != nil {
			log.Error().Err(err).Msg("metrics server failed")
		}
	}()

	// pollMarketsLoop runs an initial poll then every 15 minutes until shutdown.
	go func() {
		pollMarkets(redisClient)
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			pollMarkets(redisClient)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down market-poller...")
}

// pollMarkets fetches external market feeds, computes signals via pkg/market, and writes
// cache.KeyMarketSignals to Redis with a one-hour TTL for downstream consumers.
func pollMarkets(redis *cache.Client) {
	log.Info().Msg("polling market data (NSE, MCX, CCIL, AMFI)...")

	_ = market.FetchNSEData()
	_ = market.FetchMCXData()
	_ = market.FetchCCILData()
	_ = market.FetchAMFIData()

	signals, err := market.ComputeOpportunitySignals(nil, nil, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to compute opportunity signals")
		return
	}

	ctx := context.Background()
	if err := redis.SetJSON(ctx, cache.KeyMarketSignals, signals, time.Hour); err != nil {
		log.Error().Err(err).Msg("failed to cache market signals")
		return
	}
	log.Info().Int("signals", len(signals)).Msg("market signals cached in Redis")
}
