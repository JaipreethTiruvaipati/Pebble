//go:build integration

// Package integration holds end-to-end tests for the async bill scoring pipeline.
package integration

import (
	"testing"
)

// TestTransactionScoringPipeline verifies bill upload → bills.uploaded → scoring → bills.scored (planned).
func TestTransactionScoringPipeline(t *testing.T) {
	t.Log("Testing Bill Upload -> RabbitMQ -> Scoring Service -> DB Update flow...")
	// TODO: Phase 2 - Spin up RabbitMQ testcontainer, mock S3, verify DLQ logic
	// This ensures our async pipeline doesn't lose messages under load.
}
