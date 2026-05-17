package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("starting bill-service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	uploader := NewS3Uploader(cfg.AWSS3Bucket, cfg.AWSRegion)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// In a real environment, this route would be called by the api-gateway
	// or sit behind the gateway. We set it up here to receive the multipart form.
	r.Post("/upload", HandleBillUpload(uploader))

	log.Info().Msg("bill-service listening on port 8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
