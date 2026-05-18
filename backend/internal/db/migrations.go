// Package db provides PostgreSQL connectivity and schema migration helpers for Pebble services.
package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jaipreeth/pebble/backend/internal/config"
	"github.com/rs/zerolog/log"
)

// RunMigrations applies all pending UP SQL migrations from the Pebble migrations directory.
//
// Parameters:
//   - cfg: application config; DatabaseURL targets the postgres instance to migrate
//   - migrationsPath: absolute or relative filesystem path to backend/migrations (without file://)
//
// Returns:
//   - nil when migrations succeed or the schema is already current (migrate.ErrNoChange)
//   - error when migrate initialization or Up() fails for any other reason
//
// How it works: builds a golang-migrate instance with source file://{migrationsPath} and
// applies every unapplied *.up.sql in version order (users, wallets, transactions, penalties,
// pool_contributions, investments, referral_codes, streak indexes, etc.). Typically invoked
// from api-gateway or Makefile on deploy so query code and schema stay aligned before
// services accept traffic.
func RunMigrations(cfg *config.Config, migrationsPath string) error {
	log.Info().Str("path", migrationsPath).Msg("running database migrations...")

	// Create a new migrate instance
	// We use the file:// source to point to our local migrations directory
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	
	m, err := migrate.New(sourceURL, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}
	defer m.Close()

	// Run the migrations (Up applies all pending up.sql files)
	if err := m.Up(); err != nil {
		// migrate.ErrNoChange is returned when the DB is already fully migrated.
		// This is not an error condition for us, it just means "nothing to do".
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("database is already up to date (no migrations applied)")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info().Msg("database migrations applied successfully")
	return nil
}
