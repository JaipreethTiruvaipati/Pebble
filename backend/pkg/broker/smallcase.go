package broker

import "github.com/rs/zerolog/log"

// SmallcaseClient handles integration with the Smallcase Partner API.
type SmallcaseClient struct {
	APIKey  string
	BaseURL string
}

// NewSmallcaseClient initializes the broker client, defaulting to the sandbox environment.
func NewSmallcaseClient(apiKey string) *SmallcaseClient {
	return &SmallcaseClient{
		APIKey:  apiKey,
		BaseURL: "https://sandbox.smallcase.com/gateway", // Enforced sandbox first
	}
}

// ExecuteTrade submits an order to Smallcase to buy a specific asset class or basket.
func (c *SmallcaseClient) ExecuteTrade(assetClass string, amount float64) error {
	log.Info().
		Str("asset", assetClass).
		Float64("amount", amount).
		Msg("executing trade via Smallcase Partner API (sandbox)")
		
	// TODO: Phase 2 - Implement HTTP POST to Smallcase API with JWT/HMAC auth
	return nil
}
