// Package helpers provides shared test utilities for Pebble integration tests.
// testdb will spin up ephemeral PostgreSQL (testcontainers) and run migrations before HTTP tests.
package helpers

// TODO: implement testdb — ConnectTestDB(ctx) (*pgxpool.Pool, func(), error) for integration suite.
