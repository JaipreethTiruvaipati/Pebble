// Package broker abstracts trade execution against Pebble's partner broker (Smallcase).
//
// smallcase.go provides the production-oriented client used by investment-service
// to submit ETF and liquid-fund orders during pooled micro-batch execution.
package broker

import "github.com/rs/zerolog/log"

// SmallcaseClient holds credentials and base URL for the Smallcase Partner API.
// investment-service constructs one at startup from config and injects it into PoolExecutor.
type SmallcaseClient struct {
	APIKey  string
	BaseURL string
}

// NewSmallcaseClient builds a client for the Smallcase gateway.
//
// BaseURL defaults to the sandbox gateway so development and CI never hit live markets.
// APIKey is the partner key from config (investment-service main).
func NewSmallcaseClient(apiKey string) *SmallcaseClient {
	return &SmallcaseClient{
		APIKey:  apiKey,
		BaseURL: "https://sandbox.smallcase.com/gateway", // Enforced sandbox first
	}
}

// ExecuteTrade submits a buy order for instrument (e.g. NIFTY50_ETF) for amount INR.
//
// Unknown instrument codes fall back to LIQUID_FUND after IsValidInstrument check.
// Phase 2 will POST to Smallcase with JWT/HMAC auth; today it logs and returns nil so
// PoolExecutor can complete the DB and queue flow in sandbox mode.
func (c *SmallcaseClient) ExecuteTrade(instrument string, amount float64) error {
	if !IsValidInstrument(instrument) {
		instrument = "LIQUID_FUND"
	}
	log.Info().
		Str("instrument", instrument).
		Float64("amount", amount).
		Str("env", c.BaseURL).
		Msg("executing trade via Smallcase Partner API (sandbox)")

	// TODO: Phase 2 - Implement HTTP POST to Smallcase API with JWT/HMAC auth
	return nil
}
