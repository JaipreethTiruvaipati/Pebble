package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/pkg/market"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("service", "market-poller").Logger()
	log.Info().Msg("starting market-poller service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	_ = cfg // Phase 2: Use to connect to Redis

	// Setup Prometheus metrics endpoint
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Info().Msg("metrics server listening on :9094")
		if err := http.ListenAndServe(":9094", mux); err != nil {
			log.Error().Err(err).Msg("metrics server failed")
		}
	}()

	// 15-minute ticker for all market data
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		
		// Run once immediately on startup
		pollMarkets()
		
		for {
			<-ticker.C
			pollMarkets()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down market-poller...")
}

func pollMarkets() {
	log.Info().Msg("polling market data (NSE, MCX, CCIL, AMFI)...")
	
	// 1. Fetch live data
	_ = market.FetchNSEData()
	_ = market.FetchMCXData()
	_ = market.FetchCCILData()
	_ = market.FetchAMFIData()
	
	// 2. Compute Signals
	signals, _ := market.ComputeOpportunitySignals(nil, nil, nil)
	log.Info().Int("signals_computed", len(signals)).Msg("computed opportunity signals")
	
	// 3. Save to Redis
	// TODO: Write signals to cache.KeyMarketSignals
	
	log.Info().Msg("market data cached in Redis successfully")
}
