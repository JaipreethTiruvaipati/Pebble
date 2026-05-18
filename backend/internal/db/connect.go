// Package db provides PostgreSQL connectivity and schema migration helpers for Pebble services.
// Every backend microservice (api-gateway, bill-service, penalty-service, etc.) calls Connect
// at startup to obtain a shared pgxpool.Pool used by internal/db/queries.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/rs/zerolog/log"
)

// Connect establishes a verified connection pool to the Pebble PostgreSQL database.
//
// Parameters:
//   - ctx: cancellation context for pool creation and the initial health check
//   - cfg: application config; DatabaseURL must be a valid postgres DSN
//
// Returns:
//   - *pgxpool.Pool: a warmed pool (MinConns pre-allocated) ready for queries package calls
//   - error: parse, dial, or ping failures; the pool is closed on ping failure
//
// How it works: parses cfg.DatabaseURL into pgxpool.Config, sets production-oriented pool
// limits (25 max, 5 min, 1h lifetime), creates the pool, and Ping's before returning so
// callers fail fast at boot rather than on the first user request. Pebble uses pgxpool
// instead of database/sql for JSONB support and lower allocation overhead on hot paths
// (transactions, penalties, pool_contributions).
func Connect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	// Parse the connection string into a configuration object
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool settings
	// In production, you'd want to tune these based on your instance size
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	log.Info().Msg("connecting to database...")

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify the connection is actually alive
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("database connection established successfully")
	return pool, nil
}
