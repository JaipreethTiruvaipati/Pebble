package httputil

import "testing"

func TestValidateReferralCode(t *testing.T) {
	code, err := ValidateReferralCode("  pebble-abc123  ")
	if err != nil {
		t.Fatal(err)
	}
	if code != "PEBBLE-ABC123" {
		t.Fatalf("got %q", code)
	}
}

func TestValidateAmount(t *testing.T) {
	if err := ValidateAmount(0, 0); err == nil {
		t.Fatal("expected error for zero amount")
	}
	if err := ValidateAmount(100, 50); err == nil {
		t.Fatal("expected error for amount over max")
	}
}
