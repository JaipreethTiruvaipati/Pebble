// Package main runs the scoring-service microservice: consumes bills.uploaded from RabbitMQ,
// runs OCR (Vision) and LLM scoring (Gemini), persists line items, evaluates weekly streaks,
// and publishes bills.scored for penalty-service. Also exposes Prometheus metrics on :9091.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/db"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/models"
	"github.com/jaipreeth/pebble/backend/internal/queue"
	"github.com/jaipreeth/pebble/backend/pkg/llm"
	"github.com/jaipreeth/pebble/backend/pkg/ocr"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// BillUploadedEvent is the JSON payload on routing key bills.uploaded from bill-service.
type BillUploadedEvent struct {
	TransactionID string `json:"transaction_id"`
	S3Key         string `json:"s3_key"`
	UserID        string `json:"user_id"`
}

// main connects to PostgreSQL and RabbitMQ, registers consumer queue scoring.bills.uploaded
// on topic bills.uploaded, and blocks until SIGINT/SIGTERM.
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("service", "scoring-service").Logger()

	log.Info().Msg("starting scoring-service...")

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":9091", mux)
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx := context.Background()
	dbPool, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer dbPool.Close()

	gemini, err := llm.NewGeminiClient(ctx, cfg.GeminiAPIKey)
	if err != nil {
		log.Warn().Err(err).Msg("gemini client unavailable — scoring will use OCR stub only")
	}
	if gemini != nil {
		defer gemini.Close()
	}
	vision := ocr.NewVisionClient(cfg.GoogleVisionCredPath)

	rmq, err := queue.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to rabbitmq")
	}
	defer rmq.Close()

	// consumeBillsUploaded handles bills.uploaded messages and delegates to processBillUploaded.
	err = rmq.Consume("scoring.bills.uploaded", queue.TopicBillsUploaded, func(body []byte) error {
		var event BillUploadedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		return processBillUploaded(ctx, dbPool, rmq, gemini, vision, event)
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start consumer")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down scoring-service...")
}

// processBillUploaded OCRs the receipt, scores line items via Gemini (or dev stub), marks the
// transaction scored, runs EvaluateWeeklyStreak, and publishes bills.scored for penalty-service.
func processBillUploaded(ctx context.Context, dbPool *pgxpool.Pool, rmq *queue.RabbitMQ, gemini *llm.GeminiClient, vision *ocr.VisionClient, event BillUploadedEvent) error {
	txID, err := uuid.Parse(event.TransactionID)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return err
	}

	rawText, err := vision.ExtractText(ctx, nil)
	if err != nil {
		return err
	}

	var items []struct {
		Name, Category, Reasoning string
		Amount                    float64
		Score                     int
	}
	if gemini != nil {
		extraction, err := gemini.ExtractAndScore(ctx, rawText)
		if err != nil {
			return err
		}
		for _, it := range extraction.Items {
			items = append(items, struct {
				Name, Category, Reasoning string
				Amount                    float64
				Score                     int
			}{it.Name, it.Category, it.Reasoning, it.Amount, it.Score})
		}
	} else {
		items = append(items, struct {
			Name, Category, Reasoning string
			Amount                    float64
			Score                     int
		}{"Stub Item", "food", "dev fallback", 299, 72})
	}

	for _, it := range items {
		_, err := queries.InsertLineItem(ctx, dbPool, txID, models.ScoredItem{
			Name: it.Name, Amount: it.Amount, Score: it.Score,
			Category: it.Category, Reasoning: it.Reasoning,
		})
		if err != nil {
			return err
		}
	}
	if err := queries.MarkTransactionScored(ctx, dbPool, txID); err != nil {
		return err
	}

	if err := EvaluateWeeklyStreak(ctx, dbPool, rmq, userID); err != nil {
		log.Warn().Err(err).Msg("streak evaluation failed")
	}

	return rmq.Publish(ctx, queue.TopicBillsScored, queue.BillsScoredEvent{
		TransactionID: txID,
		UserID:        userID,
	})
}
