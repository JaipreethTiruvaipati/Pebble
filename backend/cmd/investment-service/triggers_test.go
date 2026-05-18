package main

import (
	"testing"
	"time"

	"github.com/jaipreeth/pebble/backend/internal/models"
)

func TestShouldRunMonthlySIP(t *testing.T) {
	ist := time.FixedZone("IST", 5*3600+1800)
	firstNine := time.Date(2026, 5, 1, 9, 15, 0, 0, ist)
	if !ShouldRunMonthlySIP(firstNine, ist) {
		t.Fatal("expected SIP on 1st at 9 AM IST")
	}
	secondTen := time.Date(2026, 5, 2, 9, 0, 0, 0, ist)
	if ShouldRunMonthlySIP(secondTen, ist) {
		t.Fatal("should not run on 2nd of month")
	}
}

func TestHasStrongBuySignal(t *testing.T) {
	signals := []models.MarketSignal{
		{AssetClass: "equity", Action: "BUY", Value: 85},
	}
	if !HasStrongBuySignal(signals) {
		t.Fatal("expected strong buy")
	}
	weak := []models.MarketSignal{{AssetClass: "equity", Action: "BUY", Value: 40}}
	if HasStrongBuySignal(weak) {
		t.Fatal("expected no strong buy for low value")
	}
}
