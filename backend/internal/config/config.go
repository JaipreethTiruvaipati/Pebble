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
	AppEnv string // "development" | "production"
	Port   string // HTTP listen port (default: "8080")

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

// Load reads environment variables (from .env.local in dev, from the real environment in prod).
// Required variables cause a panic — the service cannot start without them.
// Optional variables fall back to sensible defaults.
func Load() (*Config, error) {
	// In development, load from .env.local. In production (ECS), vars come from Secrets Manager
	// injected into the task environment — godotenv.Load will silently fail if file is missing.
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		// Server
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("PORT", "8080"),

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

// IsDevelopment returns true when running locally.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// ── helpers ──────────────────────────────────────────────────────────────────

// mustGetEnv panics if a required env var is missing.
// This is intentional — a misconfigured service should fail loudly at startup,
// not silently misbehave at runtime.
func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set — check .env.local or Secrets Manager", key))
	}
	return v
}

// getEnv returns the env var value or a fallback default.
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// getEnvInt parses an integer env var with a default.
func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
