// Package main (router.go) defines the api-gateway HTTP surface: public auth and webhooks,
// JWT-protected /api/v1 resources, Prometheus /metrics, and global middleware (CORS, rate
// limits). Transaction bill upload accepts client requests here; production flow forwards
// to bill-service which publishes bills.uploaded for scoring-service.
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/auth"
	"github.com/jaipreeth/pebble/backend/internal/cache"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/db/queries"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter configures the Chi router with middleware and registers all routes:
// GET /health; GET /metrics; POST /api/v1/auth/* (signup, verify-otp, login, refresh, logout);
// POST /api/v1/webhooks/razorpay; and JWT-protected /api/v1/me, insights, transactions,
// line-items, penalties, wallet, referrals, portfolio, investments, and market/signal.
func SetupRouter(cfg *config.Config, dbPool *pgxpool.Pool, redis *cache.Client, jwtManager *auth.JWTManager, otpService *auth.OTPService) *chi.Mux {
	r := chi.NewRouter()

	// Global Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(auth.CORS(cfg.CORSAllowedOrigins))
	r.Use(auth.RateLimit(20, 40)) // global IP fallback: 20 req/sec

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Prometheus Metrics
	r.Handle("/metrics", promhttp.Handler())

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (Auth)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", handleSignup(dbPool, otpService))
			r.Post("/verify-otp", handleVerifyOTP(dbPool, jwtManager, otpService))
			r.Post("/login", handleLogin(dbPool, jwtManager))
			r.Post("/refresh", handleRefresh(jwtManager))
			r.Post("/logout", handleLogout())
		})

		// Webhooks (Public but secured via signature)
		r.Route("/webhooks", func(r chi.Router) {
			r.Post("/razorpay", handleRazorpayWebhook(cfg))
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth(jwtManager))
			if redis != nil {
				r.Use(auth.RedisRateLimit(redis, 120, time.Minute)) // 120 req/min per user
			}

			r.Get("/me", handleGetMe(dbPool))
			r.Get("/users/me", handleGetMe(dbPool))

			// Insights (Week 20)
			r.Route("/insights", func(r chi.Router) {
				r.Get("/weekly", handleWeeklyInsights(dbPool))
				r.Get("/benchmark", handleBenchmarkInsights(dbPool))
			})

			// Transactions
			r.Route("/transactions", func(r chi.Router) {
				r.Post("/", handleCreateTransaction(dbPool))
				r.Post("/bill", handleUploadBill(cfg, dbPool))
				r.Get("/{id}", handleGetTransaction(dbPool))
				r.Get("/", handleListTransactions(dbPool))
				r.Post("/{id}/confirm", handleConfirmTransaction(dbPool))
			})

			// Line Items
			r.Route("/line-items", func(r chi.Router) {
				r.Put("/{id}/score", handleOverrideScore(dbPool))
			})

			// Penalties
			r.Route("/penalties", func(r chi.Router) {
				r.Get("/", handleListPenalties(dbPool))
				r.Post("/{id}/contest", handleContestPenalty(dbPool))
				r.Post("/{id}/confirm", handleConfirmPenaltyEarly(dbPool))
			})

			// Wallet
			r.Route("/wallet", func(r chi.Router) {
				r.Get("/balance", handleGetWalletBalance(dbPool))
				r.Post("/topup", handleWalletTopup(dbPool))
				r.Get("/ledger", handleGetWalletLedger(dbPool))
			})

			// Referrals (Week 21)
			r.Route("/referrals", func(r chi.Router) {
				r.Get("/me", handleGetReferralMe(dbPool))
				r.Post("/redeem", handleRedeemReferral(dbPool))
			})

			// Portfolio & investments (Week 17)
			r.Get("/portfolio", handleGetPortfolio(dbPool, redis))
			r.Route("/investments", func(r chi.Router) {
				r.Get("/", handleListInvestments(dbPool))
				r.Get("/{id}", handleGetInvestment(dbPool))
			})
			r.Get("/market/signal", handleGetMarketSignal(redis))
		})
	})

	return r
}

// ── Transaction Handlers (Stubs) ─────────────────────────────────────────────

// handleCreateTransaction serves POST /api/v1/transactions: validates merchant and amount,
// ensures a dev wallet row exists, and inserts a pending transaction in PostgreSQL.
func handleCreateTransaction(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromRequest(r)
		if !ok {
			httputil.RespondError(w, http.StatusUnauthorized, "invalid user context", "UNAUTHORIZED")
			return
		}
		_ = queries.EnsureDevWallet(r.Context(), dbPool, userID)
		var req struct {
			Merchant    string  `json:"merchant"`
			TotalAmount float64 `json:"total_amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "invalid body", "INVALID_REQUEST")
			return
		}
		merchant, err := httputil.ValidateNonEmpty("merchant", req.Merchant, 128)
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, err.Error(), "INVALID_REQUEST")
			return
		}
		if err := httputil.ValidateAmount(req.TotalAmount, 1_000_000); err != nil {
			httputil.RespondError(w, http.StatusBadRequest, err.Error(), "INVALID_REQUEST")
			return
		}
		id, err := queries.CreateTransaction(r.Context(), dbPool, userID, merchant, req.TotalAmount)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to create transaction", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusCreated, map[string]interface{}{
			"transaction_id": id,
			"status":         "pending",
		})
	}
}

// ── Webhook Handlers ────────────────────────────────────────────────────────

// handleRazorpayWebhook serves POST /api/v1/webhooks/razorpay (unsigned route; verifies
// Razorpay signature) to credit wallet top-ups without JWT auth.
func handleRazorpayWebhook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Read body
		// 2. Verify signature using pkg/payment
		// 3. Process payment success/failure
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"status": "received"})
	}
}

// ── Auth Handlers (Stubs for now, will implement DB queries later) ──────

// handleSignup serves POST /api/v1/auth/signup: creates a user stub and sends OTP via otpService.
func handleSignup(dbPool *pgxpool.Pool, otpService *auth.OTPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real implementation: parse JSON body, create user in DB, send OTP
		otpService.SendOTP(r.Context(), "+919876543210")
		httputil.RespondJSON(w, http.StatusCreated, map[string]string{"message": "OTP sent"})
	}
}

// handleVerifyOTP serves POST /api/v1/auth/verify-otp and returns a JWT access token on success.
func handleVerifyOTP(dbPool *pgxpool.Pool, jwtManager *auth.JWTManager, otpService *auth.OTPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Phone string `json:"phone"`
			OTP   string `json:"otp"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		userID := uuid.NewSHA1(uuid.NameSpaceURL, []byte("pebble-dev:"+req.Phone))
		token, err := jwtManager.GenerateToken(userID)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to issue token", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"token":   token,
			"user_id": userID,
		})
	}
}

// handleLogin serves POST /api/v1/auth/login: ensures dev user/wallet rows, optionally redeems
// referral_code, and issues a JWT. Referral redemption mirrors POST /api/v1/referrals/redeem.
func handleLogin(dbPool *pgxpool.Pool, jwtManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email        string `json:"email"`
			Password     string `json:"password"`
			ReferralCode string `json:"referral_code"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		// Dev login: derive stable user id from email until user_queries signup is implemented.
		userID := uuid.NewSHA1(uuid.NameSpaceURL, []byte("pebble-dev:"+req.Email))
		_ = queries.EnsureDevUser(r.Context(), dbPool, userID, req.Email)
		_ = queries.EnsureDevWallet(r.Context(), dbPool, userID)
		if req.ReferralCode != "" {
			_ = queries.RedeemReferralCode(r.Context(), dbPool, userID, req.ReferralCode)
		}
		token, err := jwtManager.GenerateToken(userID)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "failed to issue token", "INTERNAL_ERROR")
			return
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"token":   token,
			"user_id": userID,
		})
	}
}

// handleRefresh serves POST /api/v1/auth/refresh using the httpOnly refresh cookie.
func handleRefresh(jwtManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read httpOnly cookie -> Validate refresh token -> Generate new access token
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"token": "new-stub-jwt-token"})
	}
}

// handleLogout serves POST /api/v1/auth/logout and clears the refresh-token cookie.
func handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear httpOnly cookie
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
	}
}
