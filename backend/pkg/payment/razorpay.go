// Package payment implements payment-provider integrations for Pebble wallet top-ups.
//
// api-gateway Razorpay webhooks must verify payload authenticity before crediting balances.
// Verification uses the shared secret from config (RazorpayKeySecret).
package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifyRazorpaySignature checks that body was signed by Razorpay with secret.
//
// It computes HMAC-SHA256(hex) over body and compares to signature using hmac.Equal
// to resist timing attacks. handleRazorpayWebhook in api-gateway should call this
// before processing payment.captured or payment.failed events.
func VerifyRazorpaySignature(body []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	// Use hmac.Equal to prevent timing attacks
	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}
