package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/rs/zerolog/log"
)

// Connect establishes a connection pool to the PostgreSQL database.
// It verifies the connection with a Ping before returning.
// We use pgxpool instead of standard database/sql because pgx is significantly
// faster and provides native support for PostgreSQL-specific features (like JSONB).
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
