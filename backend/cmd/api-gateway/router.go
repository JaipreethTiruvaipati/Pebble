package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/auth"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter configures the Chi router and all API routes.
func SetupRouter(cfg *config.Config, dbPool *pgxpool.Pool, jwtManager *auth.JWTManager, otpService *auth.OTPService) *chi.Mux {
	r := chi.NewRouter()

	// Global Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(auth.RateLimit(10, 20)) // 10 req/sec, burst of 20

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
			
			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				// We can safely extract the UserID here because of RequireAuth
				userID := r.Context().Value(auth.UserIDKey)
				httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{"user_id": userID})
			})

			// Transactions
			r.Route("/transactions", func(r chi.Router) {
				r.Post("/", handleCreateTransaction(dbPool))
				r.Post("/bill", handleUploadBill(dbPool))
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
		})
	})

	return r
}

// ── Transaction Handlers (Stubs) ─────────────────────────────────────────────

func handleCreateTransaction(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.RespondJSON(w, http.StatusCreated, map[string]string{"message": "transaction created manually"})
	}
}

func handleUploadBill(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Acts as a proxy to the bill-service, or simply puts the job on the queue
		httputil.RespondJSON(w, http.StatusAccepted, map[string]string{"message": "bill uploaded and queued for processing"})
	}
}

func handleGetTransaction(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"id": id, "status": "scored"})
	}
}

func handleListTransactions(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.RespondJSON(w, http.StatusOK, []map[string]interface{}{})
	}
}

func handleConfirmTransaction(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"transaction_id": id,
			"penalties_created": 2,
			"total_penalty_queued": 145.50,
			"status": "confirmed",
		})
	}
}

func handleOverrideScore(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"line_item_id": id,
			"message": "score successfully overridden",
			"user_overridden": true,
		})
	}
}

// ── Penalty Handlers (Stubs) ────────────────────────────────────────────────

func handleListPenalties(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.RespondJSON(w, http.StatusOK, []map[string]interface{}{})
	}
}

func handleContestPenalty(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"id": id, "status": "contested"})
	}
}

func handleConfirmPenaltyEarly(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"id": id, "status": "confirmed"})
	}
}

// ── Wallet Handlers (Stubs) ─────────────────────────────────────────────────

func handleGetWalletBalance(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.RespondJSON(w, http.StatusOK, map[string]float64{
			"balance":        1500.50,
			"pending_total":  145.50,
			"invested_total": 5000.00,
		})
	}
}

func handleWalletTopup(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Expects Razorpay payment token in body
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "topup successful"})
	}
}

func handleGetWalletLedger(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.RespondJSON(w, http.StatusOK, []map[string]interface{}{})
	}
}

// ── Webhook Handlers ────────────────────────────────────────────────────────

func handleRazorpayWebhook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Read body
		// 2. Verify signature using pkg/payment
		// 3. Process payment success/failure
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"status": "received"})
	}
}

// ── Auth Handlers (Stubs for now, will implement DB queries later) ──────

func handleSignup(dbPool *pgxpool.Pool, otpService *auth.OTPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real implementation: parse JSON body, create user in DB, send OTP
		otpService.SendOTP(r.Context(), "+919876543210")
		httputil.RespondJSON(w, http.StatusCreated, map[string]string{"message": "OTP sent"})
	}
}

func handleVerifyOTP(dbPool *pgxpool.Pool, jwtManager *auth.JWTManager, otpService *auth.OTPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check OTP -> Generate JWT -> Return UserProfile
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"token": "stub-jwt-token"})
	}
}

func handleLogin(dbPool *pgxpool.Pool, jwtManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check email/password -> Generate JWT -> Return UserProfile
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"token": "stub-jwt-token"})
	}
}

func handleRefresh(jwtManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read httpOnly cookie -> Validate refresh token -> Generate new access token
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"token": "new-stub-jwt-token"})
	}
}

func handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear httpOnly cookie
		httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
	}
}
