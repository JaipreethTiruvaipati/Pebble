// Package config loads and holds environment-backed settings for all Pebble backend
// services. Load runs once at process start; services pass *Config instead of calling
// os.Getenv so secrets, TTLs, and integration endpoints stay centralized and testable.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds every environment variable the backend services need.
// It is loaded once at startup and passed around — never call os.Getenv directly in services.
type Config struct {
	// ── Server ────────────────────────────────────────────────────────────────
	AppEnv            string // "development" | "production"
	Port              string // HTTP listen port (default: "8080")
	CORSAllowedOrigins string // comma-separated browser origins

	// ── Database ──────────────────────────────────────────────────────────────
	DatabaseURL string // full PostgreSQL DSN: postgres://user:pass@host:port/db

	// ── Redis ─────────────────────────────────────────────────────────────────
	RedisURL string // redis://host:6379

	// ── RabbitMQ ─────────────────────────────────────────────────────────────
	RabbitMQURL string // amqp://user:pass@host:5672/

	// ── JWT (RS256 asymmetric signing) ────────────────────────────────────────
	JWTPrivateKeyPath    string        // path to PEM file — used by api-gateway to SIGN tokens
	JWTPublicKeyPath     string        // path to PEM file — used by all services to VERIFY tokens
	JWTAccessExpiry      time.Duration // how long an access token lives (default: 15 min)
	JWTRefreshExpiryDays int           // how long a refresh token lives (default: 30 days)

	// ── AWS ───────────────────────────────────────────────────────────────────
	AWSRegion   string // "ap-south-1" (Mumbai — RBI data localisation requirement)
	AWSS3Bucket string // bucket name for bill images / PDFs

	// ── External APIs ─────────────────────────────────────────────────────────
	GoogleVisionCredPath string // path to Google service account JSON
	GeminiAPIKey         string // Gemini API key for LLM impulse scoring

	// ── Payments ──────────────────────────────────────────────────────────────
	RazorpayKeyID     string
	RazorpayKeySecret string

	// ── Notifications ─────────────────────────────────────────────────────────
	FirebaseCredPath string // path to Firebase Admin SDK JSON (FCM push)
	SESFromEmail     string // sender address for AWS SES emails

	// ── Broker (investment execution) ─────────────────────────────────────────
	SmallcaseAPIKey    string
	SmallcaseAPISecret string
	SmallcaseEnv       string // "sandbox" | "production"

	// ── OTP ───────────────────────────────────────────────────────────────────
	OTPProvider string // "sns" | "msg91"
	Msg91AuthKey string
}

// Load reads environment variables and populates Config with required and optional values.
//
// Inputs: process environment; in development also .env.local via godotenv (ignored if missing).
// Outputs: *Config and nil error on success; panic via mustGetEnv if DATABASE_URL, REDIS_URL,
// or RABBITMQ_URL are unset.
//
// Every service main calls Load first, then passes cfg to auth, cache, db, and queue setup.
// JWT paths, CORS origins, and broker keys flow into api-gateway; worker services use subsets.
func Load() (*Config, error) {
	// In development, load from .env.local. In production (ECS), vars come from Secrets Manager
	// injected into the task environment — godotenv.Load will silently fail if file is missing.
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		// Server
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000"),

		// Database — required; the service is useless without it
		DatabaseURL: mustGetEnv("DATABASE_URL"),

		// Redis — required for rate limiting and caching
		RedisURL: mustGetEnv("REDIS_URL"),

		// RabbitMQ — required for async pipeline
		RabbitMQURL: mustGetEnv("RABBITMQ_URL"),

		// JWT
		JWTPrivateKeyPath:    getEnv("JWT_PRIVATE_KEY_PATH", "./keys/private.pem"),
		JWTPublicKeyPath:     getEnv("JWT_PUBLIC_KEY_PATH", "./keys/public.pem"),
		JWTAccessExpiry:      time.Duration(getEnvInt("JWT_ACCESS_EXPIRY_MINUTES", 15)) * time.Minute,
		JWTRefreshExpiryDays: getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 30),

		// AWS
		AWSRegion:   getEnv("AWS_REGION", "ap-south-1"),
		AWSS3Bucket: getEnv("AWS_S3_BUCKET", "pebble-bills-dev"),

		// External APIs
		GoogleVisionCredPath: getEnv("GOOGLE_VISION_CREDENTIALS_PATH", ""),
		GeminiAPIKey:         getEnv("GEMINI_API_KEY", ""),

		// Payments
		RazorpayKeyID:     getEnv("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret: getEnv("RAZORPAY_KEY_SECRET", ""),

		// Notifications
		FirebaseCredPath: getEnv("FIREBASE_CREDENTIALS_PATH", ""),
		SESFromEmail:     getEnv("SES_FROM_EMAIL", "noreply@pebble.in"),

		// Broker
		SmallcaseAPIKey:    getEnv("SMALLCASE_API_KEY", ""),
		SmallcaseAPISecret: getEnv("SMALLCASE_API_SECRET", ""),
		SmallcaseEnv:       getEnv("SMALLCASE_ENV", "sandbox"),

		// OTP
		OTPProvider:  getEnv("OTP_PROVIDER", "sns"),
		Msg91AuthKey: getEnv("MSG91_AUTH_KEY", ""),
	}

	return cfg, nil
}

// IsDevelopment reports whether APP_ENV is "development".
//
// Inputs: receiver Config after Load.
// Outputs: true for local dev behavior (e.g. OTP logged, relaxed defaults).
//
// Used by OTPService and handlers to skip external SMS charges and enable verbose paths.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// ── helpers ──────────────────────────────────────────────────────────────────

// mustGetEnv returns os.Getenv(key) or panics if the variable is empty.
//
// Inputs: key name of a required setting (e.g. "DATABASE_URL").
// Outputs: non-empty string value.
//
// Intentional fail-fast at startup so ECS/tasks with missing Secrets Manager bindings
// exit immediately rather than serving broken requests.
func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set — check .env.local or Secrets Manager", key))
	}
	return v
}

// getEnv returns the environment variable value or defaultVal when unset or empty.
//
// Inputs: key and fallback default for optional settings (PORT, CORS_ALLOWED_ORIGINS).
// Outputs: resolved string used in Load to build Config.
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// getEnvInt parses an integer environment variable or returns defaultVal.
//
// Inputs: key and default when missing or non-numeric.
// Outputs: int for JWT expiry minutes, refresh days, and similar numeric env vars.
func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
