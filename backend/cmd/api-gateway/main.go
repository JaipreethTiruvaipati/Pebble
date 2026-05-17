package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/auth"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure JSON logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("starting api-gateway...")

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// 2. Connect to Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbPool, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer dbPool.Close()

	// 3. Initialize Auth services
	jwtManager, err := auth.NewJWTManager(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize JWT manager")
	}
	otpService := auth.NewOTPService(cfg)

	// 4. Setup Router
	r := SetupRouter(cfg, dbPool, jwtManager, otpService)

	// 5. Start Server
	srv := NewServer(cfg.Port, r)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server...")

	ctxShutDown, cancelShutDown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutDown()
	if err := srv.Stop(ctxShutDown); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}
	log.Info().Msg("server exited properly")
}
