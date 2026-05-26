// Package auth provides authentication primitives for Pebble: JWT access tokens,
// HTTP middleware, OTP login, CORS, and distributed rate limiting. This file
// implements phone OTP generation and verification for passwordless login.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/rs/zerolog/log"
)

// OTPService generates and verifies one-time passwords for Indian mobile login.
// It reads OTPProvider and related settings from config and will integrate with
// AWS SNS or Msg91 in production; dev mode logs OTPs instead of sending SMS.
type OTPService struct {
	cfg *config.Config
}

// NewOTPService constructs an OTPService bound to application configuration.
//
// Inputs: cfg from config.Load (OTPProvider, Msg91AuthKey, IsDevelopment).
// Outputs: *OTPService with no external connections—stateless aside from cfg.
//
// Created at api-gateway startup alongside JWTManager for auth handler routes.
func NewOTPService(cfg *config.Config) *OTPService {
	return &OTPService{cfg: cfg}
}

// SendOTP generates a six-digit code for phone and dispatches it via the configured provider.
//
// Inputs: ctx for cancellation; phone E.164 or local format as stored in users.
// Outputs: the generated OTP string (for dev/testing) and error if provider fails.
//
// In development, logs OTP and skips SMS to avoid cost. Production will persist OTP
// in Redis with TTL and call SNS/Msg91; login handler compares user input via VerifyOTP.
func (s *OTPService) SendOTP(ctx context.Context, phone string) (string, error) {
	// Generate a cryptographically secure 6-digit OTP
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("failed to generate secure OTP: %w", err)
	}
	otp := fmt.Sprintf("%06d", binary.BigEndian.Uint64(b[:])%1000000)

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

// VerifyOTP checks whether providedOTP matches the expected value for phone.
//
// Inputs: ctx (reserved for Redis lookup); phone; providedOTP from client; expectedOTP from SendOTP/cache.
// Outputs: true if codes match.
//
// Stub compares strings directly; production will load expected OTP from Redis by phone
// and enforce expiry before issuing JWT via JWTManager.GenerateToken.
func (s *OTPService) VerifyOTP(ctx context.Context, phone, providedOTP, expectedOTP string) bool {
	// In a real implementation, we would pull the expectedOTP from Redis
	// For this stub, we just compare the strings directly.
	return providedOTP == expectedOTP
}
