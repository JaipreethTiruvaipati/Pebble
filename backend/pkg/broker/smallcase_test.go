package broker

import "testing"

func TestNewSmallcaseClient_UsesSandboxByDefault(t *testing.T) {
	c := NewSmallcaseClient("test-key")
	if c.BaseURL != "https://sandbox.smallcase.com/gateway" {
		t.Fatalf("expected sandbox URL, got %s", c.BaseURL)
	}
}

func TestExecuteTrade_Sandbox_NoError(t *testing.T) {
	c := NewSmallcaseClient("test-key")
	if err := c.ExecuteTrade("equity", 500.0); err != nil {
		t.Fatalf("sandbox ExecuteTrade should not fail: %v", err)
	}
}
