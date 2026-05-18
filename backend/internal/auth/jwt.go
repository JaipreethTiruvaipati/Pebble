// Package auth provides authentication primitives for Pebble: JWT access tokens,
// HTTP middleware, OTP login, CORS, and distributed rate limiting. This file
// implements RS256 JWT creation and verification used by the api-gateway after
// phone OTP login and by RequireAuth on protected routes.
package auth

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jaipreeth/pebble/backend/internal/config"
)

// JWTManager signs and verifies Pebble access tokens using RS256 asymmetric keys.
// The api-gateway holds the private key and issues tokens; all services that
// verify Bearer tokens only need the public key loaded from config.JWTPublicKeyPath.
type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	cfg        *config.Config
}

// CustomClaims is the JWT payload for Pebble access tokens.
// UserID identifies the authenticated user; RegisteredClaims carry standard
// expiry, issuer, and subject fields consumed by jwt/v5 during verification.
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTManager loads RSA PEM key pairs from disk paths in cfg and returns a
// manager ready to sign and verify tokens.
//
// Inputs: cfg with JWTPrivateKeyPath and JWTPublicKeyPath (typically from config.Load).
// Outputs: *JWTManager on success, or an error if keys cannot be read or parsed.
//
// Called once at api-gateway startup; private key must never be distributed to
// worker services—only the public key is needed for verification elsewhere.
func NewJWTManager(cfg *config.Config) (*JWTManager, error) {
	// Read Private Key
	privBytes, err := os.ReadFile(cfg.JWTPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	// Read Public Key
	pubBytes, err := os.ReadFile(cfg.JWTPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not read public key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	return &JWTManager{
		privateKey: privateKey,
		publicKey:  publicKey,
		cfg:        cfg,
	}, nil
}

// GenerateToken mints a short-lived access JWT for userID.
//
// Inputs: userID from the users table after successful OTP verification.
// Outputs: signed JWT string, or error if signing fails.
//
// Embeds UserID in CustomClaims and sets ExpiresAt from cfg.JWTAccessExpiry,
// Issuer "pebble-auth", and Subject to the user's UUID string. Used by login
// and token refresh handlers in the api-gateway before returning tokens to clients.
func (m *JWTManager) GenerateToken(userID uuid.UUID) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.cfg.JWTAccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "pebble-auth",
			Subject:   userID.String(),
		},
	}

	// RS256 requires signing with a private RSA key
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
}

// VerifyToken parses tokenStr, validates RS256 signature and expiry, and returns claims.
//
// Inputs: raw JWT from Authorization: Bearer header.
// Outputs: *CustomClaims when valid; error on bad signature, wrong alg, or expired token.
//
// Rejects non-RSA signing methods to prevent algorithm downgrade attacks. RequireAuth
// middleware calls this on every protected request and stores claims.UserID in context.
func (m *JWTManager) VerifyToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is exactly RS256 (prevents downgrade attacks)
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the public key for verification
		return m.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}
