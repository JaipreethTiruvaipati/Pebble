// Package broker abstracts trade execution against Pebble's partner broker (Smallcase).
//
// smallcase.go provides the production-oriented client used by investment-service
// to submit ETF and liquid-fund orders during pooled micro-batch execution.
package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jaipreeth/pebble/backend/pkg/retry"
	"github.com/rs/zerolog/log"
)

// SmallcaseClient holds credentials and base URL for the Smallcase Partner API.
// investment-service constructs one at startup from config and injects it into PoolExecutor.
type SmallcaseClient struct {
	APIKey    string
	APISecret string
	BaseURL   string
	Env       string // "sandbox" | "production"
	http      *http.Client
}

// NewSmallcaseClient builds a client for the Smallcase gateway.
//
// BaseURL defaults to the sandbox gateway so development and CI never hit live markets.
// APIKey is the partner key from config (investment-service main).
func NewSmallcaseClient(apiKey string) *SmallcaseClient {
	return &SmallcaseClient{
		APIKey:  apiKey,
		BaseURL: "https://sandbox.smallcase.com/gateway", // Enforced sandbox first
		Env:     "sandbox",
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewSmallcaseClientFull builds a fully configured client with secret and env toggle.
func NewSmallcaseClientFull(apiKey, apiSecret, env string) *SmallcaseClient {
	baseURL := "https://sandbox.smallcase.com/gateway"
	if env == "production" {
		baseURL = "https://gateway.smallcase.com"
	}
	return &SmallcaseClient{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseURL,
		Env:       env,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TradeRequest is the payload sent to the Smallcase Partner API.
type TradeRequest struct {
	Instrument string  `json:"instrument"`
	Amount     float64 `json:"amount"`
	OrderType  string  `json:"order_type"`
}

// ExecuteTrade submits a buy order for instrument (e.g. NIFTY50_ETF) for amount INR.
//
// Unknown instrument codes fall back to LIQUID_FUND after IsValidInstrument check.
// Uses exponential backoff retry (up to 5 attempts) for broker API resilience.
// In sandbox mode, logs the trade without making real HTTP calls.
func (c *SmallcaseClient) ExecuteTrade(instrument string, amount float64) error {
	if !IsValidInstrument(instrument) {
		instrument = "LIQUID_FUND"
	}

	log.Info().
		Str("instrument", instrument).
		Float64("amount", amount).
		Str("env", c.Env).
		Msg("executing trade via Smallcase Partner API")

	// In sandbox mode, simulate success without real API calls
	if c.Env == "sandbox" {
		log.Info().
			Str("instrument", instrument).
			Float64("amount", amount).
			Msg("sandbox trade executed (simulated)")
		return nil
	}

	// Production: POST to Smallcase API with retry
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	return retry.Do(ctx, retry.Aggressive(), "smallcase-execute-trade", func() error {
		payload := TradeRequest{
			Instrument: instrument,
			Amount:     amount,
			OrderType:  "BUY",
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal trade request: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/orders", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-sc-api-key", c.APIKey)

		resp, err := c.http.Do(req)
		if err != nil {
			return fmt.Errorf("smallcase API call failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			return fmt.Errorf("smallcase server error: %d", resp.StatusCode)
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("smallcase client error: %d (not retryable)", resp.StatusCode)
		}

		log.Info().
			Str("instrument", instrument).
			Float64("amount", amount).
			Int("status", resp.StatusCode).
			Msg("trade executed successfully")
		return nil
	})
}
