package allocate

import (
	"testing"

	"github.com/jaipreeth/pebble/backend/internal/models"
)

func TestComputeAllocation_NormalizesTo100(t *testing.T) {
	signals := []models.MarketSignal{
		{AssetClass: "equity", Action: "BUY"},
	}
	alloc := ComputeAllocation(signals)
	total := alloc.Equity + alloc.Gold + alloc.Bonds
	if total < 99.9 || total > 100.1 {
		t.Fatalf("expected ~100%%, got %.2f", total)
	}
}

func TestComputeAllocation_ClampsEquityMax(t *testing.T) {
	signals := []models.MarketSignal{
		{AssetClass: "equity", Action: "BUY"},
		{AssetClass: "equity", Action: "BUY"},
		{AssetClass: "equity", Action: "BUY"},
	}
	alloc := ComputeAllocation(signals)
	if alloc.Equity > 70.1 {
		t.Fatalf("equity should be clamped to 70%%, got %.2f", alloc.Equity)
	}
}
