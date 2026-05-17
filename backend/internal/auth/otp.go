package auth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/rs/zerolog/log"
)

// OTPService handles generating and verifying One-Time Passwords.
type OTPService struct {
	cfg *config.Config
}

// NewOTPService creates a new OTP service.
func NewOTPService(cfg *config.Config) *OTPService {
	return &OTPService{cfg: cfg}
}

// SendOTP generates a 6-digit OTP and sends it via the configured provider.
func (s *OTPService) SendOTP(ctx context.Context, phone string) (string, error) {
	// Generate a random 6-digit OTP
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := fmt.Sprintf("%06d", r.Intn(1000000))

	// In development, we don't actually want to send real SMS and incur costs.
	if s.cfg.IsDevelopment() {
		log.Info().Str("phone", phone).Str("otp", otp).Msg("DEV MODE: OTP generated but not sent to SMS provider")
		return otp, nil
	}

	// TODO: Phase 2 - Implement actual SMS integration (AWS SNS or Msg91)
	// For now, this acts as a stub.
	log.Info().Str("phone", phone).Msg("Sending OTP via provider...")
	
	return otp, nil
}

// VerifyOTP checks if the provided OTP is valid for the phone number.
func (s *OTPService) VerifyOTP(ctx context.Context, phone, providedOTP, expectedOTP string) bool {
	// In a real implementation, we would pull the expectedOTP from Redis
	// For this stub, we just compare the strings directly.
	return providedOTP == expectedOTP
}
