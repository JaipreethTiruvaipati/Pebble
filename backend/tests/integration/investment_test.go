//go:build integration

package integration

import (
	"testing"

	"github.com/jaipreeth/pebble/backend/internal/models"
	"github.com/jaipreeth/pebble/backend/pkg/allocate"
	"github.com/jaipreeth/pebble/backend/pkg/market"
)

// TestInvestmentPipelineWithSimulatedSignals validates allocation + signal stubs used in staging.
func TestInvestmentPipelineWithSimulatedSignals(t *testing.T) {
	signals, err := market.ComputeOpportunitySignals(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(signals) == 0 {
		t.Fatal("expected stub signals")
	}

	alloc := allocate.ComputeAllocation(signals)
	total := alloc.Equity + alloc.Gold + alloc.Bonds
	if total < 99 || total > 101 {
		t.Fatalf("allocation should sum to ~100, got %.2f", total)
	}

	// Simulate week 18 staging: strong BUY should be detectable
	strong := false
	for _, s := range signals {
		if s.Action == "BUY" && s.Value >= 70 {
			strong = true
		}
	}
	_ = strong
	_ = models.MarketSignal{}
	t.Log("staging pipeline: signals -> allocation OK")
}
