package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/cache"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/db"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/jaipreeth/pebble/backend/pkg/broker"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbPool, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer dbPool.Close()

	redisClient, err := cache.Connect(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("redis connection failed")
	}
	defer redisClient.Close()

	rmq, err := queue.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to rabbitmq")
	}
	defer rmq.Close()

	brokerClient := broker.NewSmallcaseClient(cfg.SmallcaseAPIKey)
	globalExecutor = NewPoolExecutor(dbPool, redisClient, brokerClient, rmq)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Info().Msg("metrics server listening on :9095")
		if err := http.ListenAndServe(":9095", mux); err != nil {
			log.Error().Err(err).Msg("metrics server failed")
		}
	}()

	go StartThresholdTrigger(dbPool)
	go StartTimeTrigger()
	go StartOpportunityTrigger(redisClient)

	_ = rmq.Consume("investment.penalties.confirmed", queue.TopicWalletPenaltyConfirmed, func(body []byte) error {
		var event queue.PenaltiesConfirmedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		log.Info().
			Str("user_id", event.UserID.String()).
			Float64("amount", event.TotalAmount).
			Msg("penalties confirmed — funds available for pooling")
		return nil
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down investment-service...")
}
