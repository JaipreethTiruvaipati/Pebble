package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/auth"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/jaipreeth/pebble/backend/internal/httputil"
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

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth(jwtManager))
			
			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				// We can safely extract the UserID here because of RequireAuth
				userID := r.Context().Value(auth.UserIDKey)
				httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{"user_id": userID})
			})

			// Further routes will go here (bills, transactions, wallet, etc.)
		})
	})

	return r
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
